package utils

import (
	"fmt"
	"strings"

	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("routes_walker")

func WalkRoutes(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	pathTemplate, err := route.GetPathTemplate()
	if err == nil {
		logger.Info("ROUTE:", pathTemplate)
	}
	pathRegexp, err := route.GetPathRegexp()
	if err == nil {
		logger.Info("Path regexp:", pathRegexp)
	}
	queriesTemplates, err := route.GetQueriesTemplates()
	if err == nil {
		logger.Info("Queries templates:", strings.Join(queriesTemplates, ","))
	}
	queriesRegexps, err := route.GetQueriesRegexp()
	if err == nil {
		logger.Info("Queries regexps:", strings.Join(queriesRegexps, ","))
	}
	methods, err := route.GetMethods()
	if err == nil {
		logger.Info("Methods:", strings.Join(methods, ","))
	}
	fmt.Println()
	return nil
}
