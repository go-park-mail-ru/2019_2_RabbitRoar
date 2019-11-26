package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/op/go-logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"time"
)

func NewMetricsMiddleware() echo.MiddlewareFunc {
	var log = logging.MustGetLogger("metrics")

	requestDurationSummary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "svoyak",
			Name:      "request_duration_s",
		},
		[]string{"method", "path"},
	)

	requestStatusCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "svoyak",
			Name:      "request_counter_status",
		},
		[]string{"path", "status"},
	)

	prometheus.MustRegister(
		requestDurationSummary,
		requestStatusCounter,
	)

	http.Handle(viper.GetString("server.metrics.path"), promhttp.Handler())
	go func() {
		log.Error(http.ListenAndServe(viper.GetString("server.metrics.address"), nil))
	}()

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
