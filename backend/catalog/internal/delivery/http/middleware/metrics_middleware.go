package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	instrumentationName = "github.com/omides248/catalog/http-middleware"
)

func MetricsMiddleware() echo.MiddlewareFunc {

	fmt.Println("0000000000000000000000")
	fmt.Println("0000000000000000000000")
	fmt.Println("0000000000000000000000")
	fmt.Println("0000000000000000000000")
	meter := otel.GetMeterProvider().Meter(instrumentationName)

	requestCounter, err := meter.Int64Counter("http.server.requests.count", metric.WithDescription("Total number of HTTP requests."))
	if err != nil {
		// Handle error, e.g., log it or panic
	}

	requestLatency, err := meter.Int64Histogram("http.server.requests.latency", metric.WithDescription("Latency of HTTP requests in milliseconds."))
	if err != nil {
		// Handle error
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			startTime := time.Now()

			// Proceed with the request
			err := next(c)

			c.Response().After(func() {
				duration := time.Since(startTime)
				attrs := []attribute.KeyValue{
					attribute.String("http.method", c.Request().Method),
					attribute.String("http.route", c.Path()),
					attribute.Int("http.status_code", c.Response().Status),
				}

				//fmt.Println("c.Request().Method", c.Request().Method)
				//fmt.Println("http.route", c.Path())
				//fmt.Println("http.status_code", c.Response().Status)
				//fmt.Println("http.latency", time.Since(startTime))

				requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
				requestLatency.Record(ctx, duration.Milliseconds(), metric.WithAttributes(attrs...))
			})

			return err
		}
	}
}
