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
	if he, ok := err.(*echo.HTTPError); ok {
		ctx.Logger().Error(he.Internal)

		if he.Code >= http.StatusInternalServerError {
			sentry.CaptureException(he)
		}
	}

	ctx.Logger().Error(errors.WithStack(err))

	sentry.CaptureException(err)

	ctx.Echo().DefaultHTTPErrorHandler(err, ctx)
}
