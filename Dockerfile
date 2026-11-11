FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh && sed -i 's/\r$//' /app/entrypoint.sh

# Public API. Prometheus metrics stay off this port.
EXPOSE 8080
# Internal Prometheus scrape port (METRICS_PORT, default 9091).
EXPOSE 9091
ENTRYPOINT ["/app/entrypoint.sh"]
