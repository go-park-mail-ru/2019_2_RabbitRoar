package http_utils

import (
	"github.com/labstack/echo/v4"
	"github.com/xeipuuv/gojsonschema"
	"strconv"
)

func GetIntParam(ctx echo.Context, defaultValue int) int {
	pageNumString := ctx.QueryParam("page")
	if pageNumString == "" {
		return defaultValue
	}
	pageNum, err := strconv.Atoi(pageNumString)
	if err != nil {
		return defaultValue
	}
	return pageNum
}

func ExtractErrors(es []gojsonschema.ResultError) []string {
	var errs []string
	for _, e := range es {
		errs = append(errs, e.Description())
	}
	return errs
}
