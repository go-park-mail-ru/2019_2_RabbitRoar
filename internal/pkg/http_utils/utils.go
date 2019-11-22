package http_utils

import (
	"github.com/labstack/echo/v4"
	"github.com/xeipuuv/gojsonschema"
	"strconv"
)

func GetIntParam(ctx echo.Context, name string, defaultValue int) int {
	param := ctx.QueryParam(name)
	if param == "" {
		return defaultValue
	}
	paramInt, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}
	return paramInt
}

func ExtractErrors(es []gojsonschema.ResultError) []string {
	var errs []string
	for _, e := range es {
		errs = append(errs, e.Description())
	}
	return errs
}
