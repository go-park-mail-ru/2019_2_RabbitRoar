package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/handlers"
	"github.com/gorilla/mux"
)

func Start() {
	topRouter := mux.NewRouter()
	topRouter.HandleFunc("/login", handlers.LoginHandler).Methods("POST", "OPTIONS")
	topRouter.HandleFunc("/signup", handlers.SignupHandler).Methods("POST", "OPTIONS")
	topRouter.Use(mux.CORSMethodMiddleware(topRouter))

	userRouter := topRouter.PathPrefix("/api").Subrouter()
	userRouter.HandleFunc("/user", handlers.GetProfileHandler).Methods("GET", "OPTIONS")
	userRouter.HandleFunc("/user", handlers.UpdateProfile).Methods("PUT", "OPTIONS")
	userRouter.HandleFunc("/user", handlers.LogoutHandler).Methods("DELETE", "OPTIONS")
	userRouter.Use(middleware.AuthMiddleware)

	topRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	log.Fatal(http.ListenAndServe(":3000", topRouter))
}
