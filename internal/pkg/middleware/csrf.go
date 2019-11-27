package middleware

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf"
	_csrfHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf/delivery/http"
	"github.com/labstack/echo/v4"
	"net/http"
)

type csrfMiddleware struct {
	jwtToken csrf.JwtToken
}

func NewCSRFMiddleware(token csrf.JwtToken) echo.MiddlewareFunc {
	mw := csrfMiddleware{
		jwtToken: token,
	}
	return mw.CSRFMiddlewareFunc
}

func (mw csrfMiddleware) CSRFMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		SessionID := ctx.Get("sessionID").(string)
		CSRF := ctx.Request().Header.Get(_csrfHttp.HeaderCSRFToken)

		if ok, err := mw.jwtToken.Check(SessionID, CSRF); !ok {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "Invalid CSRF Token",
				Internal: err,
			}
		}

		return next(ctx)
	}
}
