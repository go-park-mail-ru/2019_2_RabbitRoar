package http

import (
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

func ErrorHandler(err error, ctx echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok {
		ctx.Response().Status = he.Code
		ctx.Logger().Error(he.Internal)
	}
	ctx.Error(errors.WithStack(err))
	sentry.CaptureException(err)
	_ = ctx.NoContent(http.StatusInternalServerError)
}
