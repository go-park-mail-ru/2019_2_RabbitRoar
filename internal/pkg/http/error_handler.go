package http

import (
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

type errorView struct {
	message string
}

func ErrorHandler(err error, ctx echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok {
		err := errorView{
			message:err.Error(),
		}
		ctx.Logger().Error(he.Internal)
		if he.Code >= http.StatusInternalServerError {
			sentry.CaptureException(he)
		}
		_ = ctx.JSON(he.Code, err)
		return
	}

	ctx.Logger().Error(errors.WithStack(err))

	sentry.CaptureException(err)

	_ = ctx.NoContent(http.StatusInternalServerError)
}
