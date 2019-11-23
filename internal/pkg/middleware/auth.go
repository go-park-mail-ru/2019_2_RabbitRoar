package middleware

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"
	"net/http"

	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	sessionUseCase session.UseCase
}

func NewAuthMiddleware(sessionUseCase session.UseCase) echo.MiddlewareFunc {
	am := authMiddleware{
		sessionUseCase: sessionUseCase,
	}
	return am.AuthMiddlewareFunc
}

func (u *authMiddleware) AuthMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionID, err := ctx.Cookie("SessionID")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "No session.")
		}

		sess, err := u.sessionUseCase.GetByID(sessionID.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Session not found.")
		}

		ctx.Set("sessionID", sessionID.Value)
		ctx.Set("user", sess.User)
		return next(ctx)
	}
}
