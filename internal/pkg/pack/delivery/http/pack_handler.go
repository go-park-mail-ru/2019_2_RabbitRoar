package http

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack"
	"github.com/labstack/echo/v4"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
	"strconv"
)

type handler struct {
	packUseCase pack.UseCase
	packSchema  *gojsonschema.Schema
}

func NewPackHandler(
	e *echo.Echo,
	packUseCase pack.UseCase,
	authMiddleware echo.MiddlewareFunc,
	csrfMiddleware echo.MiddlewareFunc,
	packSchema *gojsonschema.Schema,
) {
	handler := handler{
		packUseCase: packUseCase,
		packSchema:  packSchema,
	}

	group := e.Group("/pack", authMiddleware)
	group.GET("", handler.list)
	group.POST("", csrfMiddleware(handler.create))
	group.GET("/offline", handler.offline)
	group.GET("/offline/public", handler.offlinePublic)
	group.DELETE("/:id", handler.delete)
	group.GET("/:id", handler.byID)
}

func prepareListView(packs []models.Pack) {
	for _, p := range packs {
		p.Description = ""
		p.Questions = nil
	}
	return
}

func sanitize(p *models.Pack) {
	//TODO: do it:)
}

func sanitizeSlice(p []models.Pack) {
	//TODO: do it:)
}

func extractErrors(es []gojsonschema.ResultError) []string {
	//TODO: rework error handling
	var errs []string
	for _, e := range es {
		errs = append(errs, e.Description())
	}
	return errs
}

func (h *handler) create(ctx echo.Context) error {
	var p models.Pack
	err := ctx.Bind(&p)
	if err != nil {
		return err
	}

	//TODO: overthink maybe better place is useCase coz pack structure business logic?
	loader := gojsonschema.NewGoLoader(p.Questions)
	res, err := h.packSchema.Validate(loader)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"error parsing pack",
		)
	}
	if !res.Valid() {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			extractErrors(res.Errors()),
		)
	}

	caller := ctx.Get("user").(*models.User)
	if err := h.packUseCase.Create(&p, *caller); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, p)
}

func (h *handler) offline(ctx echo.Context) error {
	caller := ctx.Get("user").(*models.User) //TODO: err on empty ?

	ids, err := h.packUseCase.FetchOffline(*caller)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, ids)
}

func (h *handler) offlinePublic(ctx echo.Context) error {
	pids, err := h.packUseCase.FetchOfflinePublic()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, pids)
}

func (h *handler) list(ctx echo.Context) error {
	//TODO: implement params
	pids, err := h.packUseCase.FetchOrderedByRating(true, 0, 10)
	if err != nil {
		return nil
	}
	prepareListView(pids)
	return ctx.JSON(http.StatusOK, pids)
}

func (h *handler) delete(ctx echo.Context) error {
	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "invalid pack id",
			Internal: err,
		}
	}
	caller := ctx.Get("user").(*models.User)
	return h.packUseCase.Delete(ID, *caller)
}

func (h *handler) byID(ctx echo.Context) error {
	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "invalid pack id",
			Internal: err,
		}
	}
	caller := ctx.Get("user").(*models.User)
	p, err := h.packUseCase.GetByID(ID, *caller)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "no pack with such ID",
			Internal: err,
		}
	}
	return ctx.JSON(http.StatusOK, p)
}
