package middleware

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	useCase user.UseCase
}

func NewAuthMiddleware(useCase user.UseCase) echo.MiddlewareFunc {
	am := authMiddleware{
		useCase: useCase,
	}
	return am.AuthMiddlewareFunc
}

func (u *authMiddleware) AuthMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionID, err := ctx.Cookie("SessionID")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "No session.")
		}

		u, err := u.useCase.GetBySessionID(sessionID.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Session not found.")
		}

		ctx.Set("sessionID", sessionID.Value)
		ctx.Set("user", u)
		return next(ctx)
	}
}
