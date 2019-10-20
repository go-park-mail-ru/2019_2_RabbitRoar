package http

import (
	"errors"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
	"github.com/labstack/echo"
)

type handler struct {
	useCase user.UseCase
}

func NewUserHandler(e *echo.Echo, usecase user.UseCase, authMiddleware echo.MiddlewareFunc) {
	handler := &handler{
		useCase: usecase,
	}

	group := e.Group("/user", authMiddleware)
	group.GET("/", handler.self)
	group.PUT("/", handler.update)
	group.PUT("/avatar", handler.avatar)
	group.GET("/:id", handler.byID)
}

func (uh *handler) self(ctx echo.Context) error {
	u := ctx.Get("user").(*models.User)
	return ctx.JSON(http.StatusOK, *u)
}

func (uh *handler) update(ctx echo.Context) error {
	return nil
}

func (uh *handler) avatar(ctx echo.Context) error {
	return nil
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

	err = ctx.JSON(http.StatusOK, *u)
	if err != nil {
		return errors.New("error marshalling user")
	}

	return nil
}
