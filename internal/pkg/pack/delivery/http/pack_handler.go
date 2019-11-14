package http

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack"
	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
	"strconv"
)

type handler struct {
	packUseCase pack.UseCase
	packSchema *gojsonschema.Schema
	sanitizer *bluemonday.Policy
}

func NewPackHandler(
	e *echo.Echo,
	packUseCase pack.UseCase,
	authMiddleware echo.MiddlewareFunc,
	csrfMiddleware echo.MiddlewareFunc,
	packSchema *gojsonschema.Schema,
) {
	handler := handler {
		packUseCase: packUseCase,
		packSchema: packSchema,
		sanitizer: bluemonday.UGCPolicy(),
	}

	group := e.Group("/pack", authMiddleware)
	group.POST("", csrfMiddleware(handler.create))
	group.DELETE("/:id", handler.delete)
	group.GET("", handler.list)
	group.GET("/offline", handler.offline)
	group.GET("/offline/public", handler.offlinePublic)
	group.GET("/author", handler.listAuthor)
	group.GET("/:id", handler.byID)
}

func (h *handler)sanitizeQuestions(p *interface{}) {
	//TODO: implement me
}

func (h *handler)sanitize(p models.Pack) models.Pack {
	p.Name = h.sanitizer.Sanitize(p.Name)
	p.Description = h.sanitizer.Sanitize(p.Description)
	p.Tags = h.sanitizer.Sanitize(p.Tags)
	h.sanitizeQuestions(&p.Questions)
	return p
}

func (h *handler)sanitizeSlice(p []models.Pack) []models.Pack {
	for i := 0; i < len(p); i++ {
		p[i] = h.sanitize(p[i])
	}
	return p
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
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "error creating pack",
			Internal: err,
		}
	}

	return ctx.JSON(http.StatusCreated, h.sanitize(p))
}

func (h *handler) offline(ctx echo.Context) error {
	caller := ctx.Get("user").(*models.User)

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
	order := ctx.Param("order")
	desc := false
	if order == "rating" {
		desc = true
	}
	var page = 0
	pageParam := ctx.Param("page")
	pageInt, err := strconv.Atoi(pageParam)
	if err != nil || pageInt < 0 {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "bad page number",
			Internal: err,
		}
	}
	page = pageInt

	packs, err := h.packUseCase.FetchOrderedByRating(desc, page, 20)
	if err != nil {
		return nil
	}
	return ctx.JSON(http.StatusOK, packs)
}

func (h* handler) listAuthor(ctx echo.Context) error {
	caller := ctx.Get("user").(*models.User)
	packs, err := h.packUseCase.FetchByAuthor(*caller, true, 0, 20)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "error fetching packs by author",
			Internal: err,
		}
	}
	return ctx.JSON(http.StatusOK, packs)
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
