package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/op/go-logging"
)

var logLog = logging.MustGetLogger("middleware_log")

func LogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		logLog.Infof("%s %s", ctx.Request().Method, ctx.Request().RequestURI)
		return next(ctx)
	}
}

