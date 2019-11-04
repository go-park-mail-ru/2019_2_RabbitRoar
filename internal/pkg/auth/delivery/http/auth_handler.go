package http

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type handler struct {
	sessionUseCase session.UseCase
	userUseCase    user.UseCase
}

func NewAuthHandler(
	e *echo.Echo,
	userUseCase user.UseCase,
	sessionUseCase session.UseCase,
	authMiddleware echo.MiddlewareFunc,
	) {

	h := handler{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}

	e.POST("/signup", h.signUp)
	e.POST("/login", h.login)
	e.DELETE("/logout", authMiddleware(h.logout))
}

func setSessionCookie(UUID *uuid.UUID, ctx echo.Context) {
	cookie := http.Cookie{
		Name:     "SessionID",
		Value:    UUID.String(),
		Expires:  time.Now().Add(512 * time.Hour),
		Secure:   false, //TODO: make me secure after ssl
		HttpOnly: true,
	}
	ctx.SetCookie(&cookie)
}

func (h *handler) signUp(ctx echo.Context) error {
	var u models.User
	err := ctx.Bind(&u)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "cannot parse user object",
			Internal: err,
		}
	}

	uc, err := h.userUseCase.Create(u)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusConflict,
			Message:  "error user with such username or password already exists",
			Internal: err,
		}
	}

	UUID, err := h.sessionUseCase.Create(*uc)
	setSessionCookie(UUID, ctx)
	return ctx.NoContent(http.StatusCreated)
}

func (h *handler) login(ctx echo.Context) error {
	var u models.User
	err := ctx.Bind(&u)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "cannot parse user object",
			Internal: err,
		}
	}

	uv, ok := h.userUseCase.IsPasswordCorrect(u)
	if !ok {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	UUID, err := h.sessionUseCase.Create(*uv)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "error creating session",
			Internal: err,
		}
	}

	setSessionCookie(UUID, ctx)

	return ctx.NoContent(http.StatusOK)
}

func (h *handler) logout(ctx echo.Context) error {
	UUID := ctx.Get("sessionID").(uuid.UUID)
	h.sessionUseCase.Destroy(UUID)
	return nil
}
