package middleware

import (
	"github.com/google/uuid"
	"github.com/op/go-logging"
	"net/http"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"

	"github.com/labstack/echo"
)

var log = logging.MustGetLogger("auth_middleware")

type authMiddleware struct {
	useCase session.UseCase
}

func NewAuthMiddleware(useCase session.UseCase) echo.MiddlewareFunc {
	am := authMiddleware{
		useCase: useCase,
	}
	return am.AuthMiddlewareFunc
}

//TODO: make session renew on close to expire
func (u *authMiddleware) AuthMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		SessionID, err := ctx.Cookie("SessionID")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "No session.")
		}

		UUID, err := uuid.Parse(SessionID.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Session invalid.")
		}

		user, err := u.useCase.GetUserByUUID(UUID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Session not found.")
		}

		ctx.Set("sessionID", UUID)
		ctx.Set("user", user)
		logger.Debugf("set sessionID %v", UUID)
		logger.Debugf("set user %v", user)
		return next(ctx)
	}
}
