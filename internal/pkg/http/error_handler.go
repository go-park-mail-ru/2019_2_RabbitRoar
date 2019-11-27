package http

import (
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

type UserErrorView struct {
	message string
}

func ErrorHandler(err error, ctx echo.Context) {
	if err == nil {
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {

		if he.Code >= http.StatusInternalServerError {
			ctx.Logger().Error(he.Internal)
			sentry.CaptureException(he)
		}

		ctx.Echo().DefaultHTTPErrorHandler(err, ctx)
		return
	}

	ctx.Logger().Error(err)

	sentry.CaptureException(errors.WithStack(err))

	ctx.Echo().DefaultHTTPErrorHandler(err, ctx)
}
