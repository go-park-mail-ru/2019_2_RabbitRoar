package middleware

import (
	"github.com/labstack/echo"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("middleware_panic")

func PanicMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			err := recover()
			logger.Critical(err)
		}()
		return next(ctx)
	}
}
