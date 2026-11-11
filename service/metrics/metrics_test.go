package metrics

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func TestRouteLabelIsTemplate(t *testing.T) {
	reg, app := newTestApp(t, "tongla-account-backend")

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		if c.Params("id") == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// Two different ids must share one series. The label is the template.
	for _, id := range []string{"42", "7"} {
		resp := mustRequest(t, app, http.MethodGet, "/users/"+id)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	}

	samples := mustCounter(t, reg, requestsMetric)
	require.Len(t, samples, 1)
	require.Equal(t, map[string]string{
		"service": "tongla-account-backend",
		"method":  "GET",
		"route":   "/users/:id",
		"status":  "200",
	}, labelMap(samples[0]))
	require.Equal(t, 2.0, samples[0].GetCounter().GetValue())
	assertRouteNotRaw(t, samples, "/users/42")
	assertRouteNotRaw(t, samples, "/users/7")

	hist := mustHistogram(t, reg, durationMetric)
	require.Len(t, hist, 1)
	require.Equal(t, "/users/:id", labelMap(hist[0])["route"])
	assertBuckets(t, hist[0])
}

func TestRouteLabelWhenGroupMiddlewareReturnsEarly(t *testing.T) {
	reg, app := newTestApp(t, "tongla-account-backend")

	// Same shape as the account and auth groups: middleware responds
	// before the registered handler is selected.
	account := app.Group("/account", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	})
	account.Get("/all-device", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	users := app.Group("/users", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	})
	users.Get("/:id", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	refresh := app.Group("/auth/refresh", func(c *fiber.Ctx) error {
		return c.Next()
	})
	refresh.Post("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	resp := mustRequest(t, app, http.MethodGet, "/account/all-device")
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp = mustRequest(t, app, http.MethodGet, "/users/42")
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp = mustRequest(t, app, http.MethodPost, "/auth/refresh")
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	samples := mustCounter(t, reg, requestsMetric)
	got := map[string]string{}
	for _, sample := range samples {
		labels := labelMap(sample)
		got[labels["method"]+" "+labels["route"]] = labels["status"]
		require.NotEqual(t, "/users/42", labels["route"])
		require.NotEqual(t, "/account", labels["route"])
		require.NotEqual(t, "/users", labels["route"])
	}
	require.Equal(t, map[string]string{
		"GET /account/all-device": "401",
		"GET /users/:id":          "401",
		"POST /auth/refresh/":     "204",
	}, got)
}

func TestUnmatchedRouteLabel(t *testing.T) {
	reg, app := newTestApp(t, "svc")
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	resp := mustRequest(t, app, http.MethodGet, "/no-such-route")
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp = mustRequest(t, app, http.MethodPost, "/users/42")
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	samples := mustCounter(t, reg, requestsMetric)
	require.Len(t, samples, 2)
	for _, sample := range samples {
		labels := labelMap(sample)
		require.Equal(t, routeUnmatched, labels["route"])
		require.Equal(t, "svc", labels["service"])
		switch labels["method"] {
		case http.MethodGet:
			require.Equal(t, "404", labels["status"])
		case http.MethodPost:
			require.Equal(t, "405", labels["status"])
		default:
			t.Fatalf("unexpected method %s", labels["method"])
		}
	}
}

func TestHandlerErrorStatus(t *testing.T) {
	reg, app := newTestApp(t, "svc")
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusTeapot, "nope")
	})
	app.Post("/explode", func(c *fiber.Ctx) error {
		return io.EOF
	})

	resp := mustRequest(t, app, http.MethodGet, "/users/1")
	require.Equal(t, http.StatusTeapot, resp.StatusCode)
	resp = mustRequest(t, app, http.MethodPost, "/explode")
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	got := map[string]string{}
	for _, sample := range mustCounter(t, reg, requestsMetric) {
		labels := labelMap(sample)
		got[labels["route"]] = labels["status"]
	}
	require.Equal(t, map[string]string{
		"/users/:id": "418",
		"/explode":   "500",
	}, got)
}

func TestSkipHealthAndMetrics(t *testing.T) {
	reg, app := newTestApp(t, "svc")
	for _, path := range []string{"/health", "/healthz", "/livez", "/readyz", "/live", "/ready", "/metrics", "/health/live"} {
		app.Get(path, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})
	}
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusCreated)
	})

	for _, path := range []string{"/health", "/health/", "/HEALTHZ", "/livez", "/readyz", "/live", "/ready", "/metrics", "/health/live"} {
		resp := mustRequest(t, app, http.MethodGet, path)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	}
	resp := mustRequest(t, app, http.MethodGet, "/users/5")
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	samples := mustCounter(t, reg, requestsMetric)
	require.Len(t, samples, 1)
	require.Equal(t, "/users/:id", labelMap(samples[0])["route"])
	require.Equal(t, "201", labelMap(samples[0])["status"])
}

func TestInFlightGauge(t *testing.T) {
	reg, app := newTestApp(t, "tongla-account-backend")
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		samples := mustGauge(t, reg, inFlightMetric)
		require.Len(t, samples, 1)
		require.Equal(t, map[string]string{"service": "tongla-account-backend"}, labelMap(samples[0]))
		require.Equal(t, 1.0, samples[0].GetGauge().GetValue())
		return c.SendStatus(fiber.StatusOK)
	})

	resp := mustRequest(t, app, http.MethodGet, "/users/1")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	samples := mustGauge(t, reg, inFlightMetric)
	require.Equal(t, 0.0, samples[0].GetGauge().GetValue())
}

func TestRuntimeAndProcessCollectors(t *testing.T) {
	reg, _ := newHTTPMetrics("tongla-account-backend")
	families, err := reg.Gather()
	require.NoError(t, err)
	require.True(t, hasMetricPrefix(families, "go_"), "expected go runtime collectors")
	require.True(t, hasMetricPrefix(families, "process_"), "expected process collectors")
}

func TestMetricsHandlerExposition(t *testing.T) {
	reg, handler := newHTTPMetrics("tongla-account-backend")
	app := fiber.New()
	app.Use(handler)
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	resp := mustRequest(t, app, http.MethodGet, "/users/42")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	rr := httptest.NewRecorder()
	metricsHandler(reg).ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	require.Equal(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	require.Contains(t, body, "http_server_requests_total")
	require.Contains(t, body, "http_server_request_duration_seconds")
	require.Contains(t, body, "http_server_requests_in_flight")
	require.Contains(t, body, `route="/users/:id"`)
	require.Contains(t, body, `service="tongla-account-backend"`)
	require.NotContains(t, body, `route="/users/42"`)
	require.Contains(t, body, "go_")
	require.Contains(t, body, "process_")
}

func TestListenPort(t *testing.T) {
	t.Run("unset defaults to 9091", func(t *testing.T) {
		unsetEnv(t, "METRICS_PORT")
		port, enabled := listenPort()
		require.True(t, enabled)
		require.Equal(t, "9091", port)
	})
	t.Run("empty disables", func(t *testing.T) {
		t.Setenv("METRICS_PORT", "")
		_, enabled := listenPort()
		require.False(t, enabled)
	})
	t.Run("zero disables", func(t *testing.T) {
		t.Setenv("METRICS_PORT", "0")
		_, enabled := listenPort()
		require.False(t, enabled)
	})
	t.Run("custom port", func(t *testing.T) {
		t.Setenv("METRICS_PORT", "19091")
		port, enabled := listenPort()
		require.True(t, enabled)
		require.Equal(t, "19091", port)
	})
}

func TestServiceName(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		unsetEnv(t, "METRICS_SERVICE_NAME")
		require.Equal(t, "tongla-account-backend", serviceName())
	})
	t.Run("override", func(t *testing.T) {
		t.Setenv("METRICS_SERVICE_NAME", "custom-svc")
		require.Equal(t, "custom-svc", serviceName())
	})
}

func TestInitServesSeparatePort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	require.NoError(t, ln.Close())

	t.Setenv("METRICS_PORT", strconv.Itoa(port))
	t.Setenv("METRICS_SERVICE_NAME", "tongla-account-backend")

	app := fiber.New()
	Init(app)
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	resp := mustRequest(t, app, http.MethodGet, "/users/3")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	url := "http://127.0.0.1:" + strconv.Itoa(port) + "/metrics"
	var body string
	require.Eventually(t, func() bool {
		res, err := http.Get(url)
		if err != nil {
			return false
		}
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil || res.StatusCode != http.StatusOK {
			return false
		}
		body = string(b)
		return strings.Contains(body, `route="/users/:id"`)
	}, 2*time.Second, 20*time.Millisecond)
	require.Contains(t, body, `status="200"`)
	require.Contains(t, body, `service="tongla-account-backend"`)
	require.NotContains(t, body, "/users/3")
}

func newTestApp(t *testing.T, service string) (*prometheus.Registry, *fiber.App) {
	t.Helper()
	reg, handler := newHTTPMetrics(service)
	app := fiber.New()
	app.Use(handler)
	return reg, app
}

func mustRequest(t *testing.T, app *fiber.App, method, path string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func mustCounter(t *testing.T, reg *prometheus.Registry, name string) []*dto.Metric {
	t.Helper()
	return mustFamily(t, reg, name)
}

func mustHistogram(t *testing.T, reg *prometheus.Registry, name string) []*dto.Metric {
	t.Helper()
	return mustFamily(t, reg, name)
}

func mustGauge(t *testing.T, reg *prometheus.Registry, name string) []*dto.Metric {
	t.Helper()
	return mustFamily(t, reg, name)
}

func mustFamily(t *testing.T, reg *prometheus.Registry, name string) []*dto.Metric {
	t.Helper()
	families, err := reg.Gather()
	require.NoError(t, err)
	for _, family := range families {
		if family.GetName() == name {
			return family.Metric
		}
	}
	t.Fatalf("metric %s not found", name)
	return nil
}

func labelMap(metric *dto.Metric) map[string]string {
	labels := make(map[string]string, len(metric.GetLabel()))
	for _, label := range metric.GetLabel() {
		labels[label.GetName()] = label.GetValue()
	}
	return labels
}

func assertRouteNotRaw(t *testing.T, metrics []*dto.Metric, raw string) {
	t.Helper()
	for _, metric := range metrics {
		require.NotEqual(t, raw, labelMap(metric)["route"])
	}
}

func assertBuckets(t *testing.T, metric *dto.Metric) {
	t.Helper()
	buckets := metric.GetHistogram().GetBucket()
	require.Len(t, buckets, len(durationBuckets))
	for i, want := range durationBuckets {
		require.InDelta(t, want, buckets[i].GetUpperBound(), 1e-12)
	}
}

func hasMetricPrefix(families []*dto.MetricFamily, prefix string) bool {
	for _, family := range families {
		if strings.HasPrefix(family.GetName(), prefix) {
			return true
		}
	}
	return false
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	prev, ok := os.LookupEnv(key)
	require.NoError(t, os.Unsetenv(key))
	t.Cleanup(func() {
		if ok {
			_ = os.Setenv(key, prev)
		} else {
			_ = os.Unsetenv(key)
		}
	})
}
