package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/op/go-logging"
	"math/rand"
	"time"
)

var logLog = logging.MustGetLogger("middleware_log")

func LogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		RID := rand.Int()
		now := time.Now().Nanosecond()
		ctx.Set("RID", RID)
		logLog.Infof(
			"RID%d pre %s %s",
			RID,
			ctx.Request().Method,
			ctx.Request().RequestURI,
		)
		err := next(ctx)
		logLog.Infof(
			"RID%d T%dns %s %s",
			RID,
			now - time.Now().Nanosecond(),
			ctx.Request().Method,
			ctx.Response().Status,
			ctx.Request().RequestURI,
		)
		return err
	}
}

