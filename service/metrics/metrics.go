// Package metrics exposes Prometheus HTTP metrics for every API route.
//
// The scrape endpoint listens on METRICS_PORT (default 9091), separate from
// the public API port, so it is not reachable through the ingress. An empty
// METRICS_PORT or "0" disables the endpoint. METRICS_SERVICE_NAME overrides
// the service label (default tongla-account-backend).
package metrics

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultServiceName = "tongla-account-backend"
	defaultMetricsPort = "9091"

	requestsMetric = "http_server_requests_total"
	durationMetric = "http_server_request_duration_seconds"
	inFlightMetric = "http_server_requests_in_flight"

	routeUnmatched = "unmatched"
)

// durationBuckets is the shared latency histogram contract, in seconds.
var durationBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30}

// Init installs the HTTP metrics middleware and, unless disabled, serves
// /metrics on METRICS_PORT. Middleware is registered before other middleware
// and routes so it observes the full request, including group middleware that
// returns before the route handler runs.
func Init(app *fiber.App) {
	port, enabled := listenPort()
	if !enabled {
		log.Printf("prometheus metrics disabled")
		return
	}

	service := serviceName()
	reg, handler := newHTTPMetrics(service)
	addr := ":" + port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("prometheus metrics listen %s: %v", addr, err)
		return
	}

	app.Use(handler)
	server := &http.Server{
		Addr:              addr,
		Handler:           metricsHandler(reg),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("prometheus metrics listening on %s service=%s", addr, service)
		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("prometheus metrics server error: %v", err)
		}
	}()
}

func serviceName() string {
	if name := os.Getenv("METRICS_SERVICE_NAME"); name != "" {
		return name
	}
	return defaultServiceName
}

// listenPort reports the metrics address. METRICS_PORT unset uses 9091.
// An empty value or "0" disables the endpoint.
func listenPort() (string, bool) {
	port, ok := os.LookupEnv("METRICS_PORT")
	if !ok {
		return defaultMetricsPort, true
	}
	port = strings.TrimSpace(port)
	if port == "" || port == "0" {
		return "", false
	}
	return port, true
}

func newHTTPMetrics(service string) (*prometheus.Registry, fiber.Handler) {
	reg := prometheus.NewRegistry()
	labels := []string{"method", "route", "status"}
	constLabels := prometheus.Labels{"service": service}

	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        requestsMetric,
		Help:        "Total number of HTTP requests processed, labeled by method, route template, and status code.",
		ConstLabels: constLabels,
	}, labels)
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        durationMetric,
		Help:        "Latency of HTTP requests in seconds, labeled by method, route template, and status code.",
		ConstLabels: constLabels,
		Buckets:     durationBuckets,
	}, labels)
	inFlight := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        inFlightMetric,
		Help:        "Number of HTTP requests currently being served.",
		ConstLabels: constLabels,
	})
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		requests,
		duration,
		inFlight,
	)

	handler := func(c *fiber.Ctx) error {
		// Do not read the request body; handlers must see it unchanged.
		if skipPath(c.Path()) {
			return c.Next()
		}

		inFlight.Inc()
		defer inFlight.Dec()

		start := time.Now()
		err := c.Next()

		// Clone method: Fiber's method string can alias the fasthttp buffer.
		method := strings.Clone(c.Method())
		route := strings.Clone(routeTemplate(c))
		status := strconv.Itoa(statusCode(c, err))
		requests.WithLabelValues(method, route, status).Inc()
		duration.WithLabelValues(method, route, status).Observe(time.Since(start).Seconds())
		return err
	}
	return reg, handler
}

func metricsHandler(reg *prometheus.Registry) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	return mux
}

// routeTemplate returns the path template registered on the router.
//
// Group middleware (authentication, for example) often responds before the
// route handler runs. Fiber then leaves the matched route pointing at that
// middleware, whose path is only the group prefix. Walking the concrete
// routes Fiber registered yields the template, such as /users/:id, rather
// than the raw path or the group prefix. Requests that match no registered
// route use "unmatched".
func routeTemplate(c *fiber.Ctx) string {
	app := c.App()
	if app == nil {
		return routeUnmatched
	}
	method := c.Method()
	path := requestPath(c)
	cfg := app.Config()
	for _, route := range app.GetRoutes(true) {
		if route.Method != method || route.Path == "" {
			continue
		}
		if fiber.RoutePatternMatch(path, route.Path, cfg) {
			return route.Path
		}
	}
	return routeUnmatched
}

// requestPath is the path Fiber used for matching, with the trailing slash
// removed when strict routing is off. The returned value is only used to
// compare against registered templates and is never stored as a label.
func requestPath(c *fiber.Ctx) string {
	path := c.Path()
	if path == "" {
		return "/"
	}
	app := c.App()
	if app != nil && !app.Config().StrictRouting && len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}
	return path
}

func statusCode(c *fiber.Ctx, err error) int {
	if err != nil {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return fiberErr.Code
		}
		// Fiber's default error handler turns every other error into 500
		// after this middleware returns, so record that status here.
		return fiber.StatusInternalServerError
	}
	code := c.Response().StatusCode()
	if code == 0 {
		return fiber.StatusOK
	}
	return code
}

// skipPath reports paths that must not be measured: the metrics scrape and
// health probes. Matching is case-insensitive and ignores a trailing slash.
func skipPath(path string) bool {
	if path == "" {
		return false
	}
	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}
	path = strings.ToLower(path)
	switch path {
	case "/metrics", "/health", "/healthz", "/livez", "/readyz", "/live", "/ready":
		return true
	default:
		return strings.HasPrefix(path, "/health/")
	}
}
