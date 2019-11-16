package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"
)

func NewMetricsMiddleware() echo.MiddlewareFunc {
	requestDurationSummary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:        "request_duration_ms",
		},
		[]string{"method", "path"},
	)

	requestStatusCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "request_counter_status",
		},
		[]string{"path", "status"},
	)

	prometheus.MustRegister(
		requestDurationSummary,
		requestStatusCounter,
	)

	http.Handle("/metrics", promhttp.Handler())
	go log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()
			err := next(ctx)
			requestDurationSummary.WithLabelValues(ctx.Request().Method, ctx.Path()).Observe(time.Since(start).Seconds())
			var status = ctx.Response().Status
			if err != nil {
				status = err.(*echo.HTTPError).Code
			}
			requestStatusCounter.WithLabelValues(ctx.Path(), strconv.Itoa(status)).Inc()
			return err
		}
	}
}
