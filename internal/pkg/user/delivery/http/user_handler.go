package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
	"github.com/labstack/echo"
)

type UserHandler struct {
	useCase user.UseCase
}

func NewUserHandler(e *echo.Group, usecase user.UseCase) {
	handler := &UserHandler{
		useCase: usecase,
	}
	e.GET("/:id", handler.UserByID)
}

func (uh *UserHandler) UserByID(ctx echo.Context) error {
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
