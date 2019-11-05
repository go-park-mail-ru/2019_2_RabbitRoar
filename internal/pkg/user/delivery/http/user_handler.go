package http

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase user.UseCase
}

func NewUserHandler(
	e *echo.Echo, usecase user.UseCase,
	authMiddleware echo.MiddlewareFunc,
	csrfMiddleware echo.MiddlewareFunc,
	) {

	handler := handler{
		useCase: usecase,
	}

	group := e.Group("/user", authMiddleware)
	group.GET("/", handler.self)
	group.PUT("/", csrfMiddleware(handler.update))
	group.PUT("/avatar", middleware.BodyLimit("2M")(csrfMiddleware(handler.avatar)))
	group.GET("/:id", handler.byID)
}

func (uh *handler) self(ctx echo.Context) error {
	u := ctx.Get("user").(*models.User)
	return ctx.JSON(http.StatusOK, uh.useCase.Sanitize(*u))
}

func (uh *handler) update(ctx echo.Context) error {
	var userUpdate models.User
	err := ctx.Bind(&userUpdate)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "can not parse user object",
			Internal: err,
		}
	}

	u := ctx.Get("user").(*models.User)

	err = uh.useCase.Update(*u, userUpdate)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "error applying user update",
			Internal: err,
		}
	}

	return ctx.JSON(http.StatusOK, uh.useCase.Sanitize(*u))
}

func (uh *handler) avatar(ctx echo.Context) error {
	file, err := ctx.FormFile("avatar")
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "error handling file from form data",
			Internal: err,
		}
	}

	u := ctx.Get("user").(*models.User)

	u, err = uh.useCase.UpdateAvatar(*u, file)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "error processing file",
			Internal: err,
		}
	}

	return ctx.JSON(http.StatusOK, uh.useCase.Sanitize(*u))
}

func (uh *handler) byID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "wrong user id provided")
	}

	u, err := uh.useCase.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "error user not found")
	}

	return ctx.JSON(http.StatusOK, uh.useCase.Sanitize(*u))
}
