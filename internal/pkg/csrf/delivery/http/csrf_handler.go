package http

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/op/go-logging"
	"net/http"
)

var log = logging.MustGetLogger("CSRF_middleware")

var HeaderCSRFToken = "X-CSRF-Token"

type handler struct {
	jwtToken csrf.JwtToken
}

type csrfToken struct {
	CSRF string
}

func NewCSRFHandler(e* echo.Echo, token csrf.JwtToken, authMiddleware echo.MiddlewareFunc) {
	h := handler{
		jwtToken:token,
	}

	e.GET("/csrf", authMiddleware(h.createCSRF))
}

func (h *handler)createCSRF(ctx echo.Context) error {
	SessionID := ctx.Get("sessionID")
	log.Debugf("got SessionID %v", SessionID)
	jwtCSRFToken, err := h.jwtToken.Create(SessionID.(uuid.UUID))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Error creating CSRF Token",
			Internal: err,
		}
	}
	return ctx.JSON(http.StatusCreated, csrfToken{CSRF:jwtCSRFToken})
}
