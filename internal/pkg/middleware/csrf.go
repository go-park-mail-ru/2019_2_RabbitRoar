package middleware

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf"
	"github.com/google/uuid"
	"github.com/labstack/echo"
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
		SessionID := ctx.Get("sessionID").(uuid.UUID)
		CSRF := ctx.Request().Header.Get("X-CSRF-Token")

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
