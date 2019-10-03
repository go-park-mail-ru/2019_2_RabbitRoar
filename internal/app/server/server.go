package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/entity"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/repository"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")
	if r.Method == http.MethodOptions {
		return
	}

	user := context.Get(r, "user").(entity.User)
	user.Password = ""
	userJSON, err := json.Marshal(user)

	if err != nil {
		panic("error marshaling user")
	}

	fmt.Println("User_profle: ", user)
	w.Write(userJSON)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	var mainUser entity.User
	mainUser = context.Get(r, "user").(entity.User)

	userUpdated := &entity.User{}

	if userFromBody, err := ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		if err := json.Unmarshal(userFromBody, userUpdated); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	changes := 0
	if userUpdated.Email != "" {
		changes++
		mainUser.Email = userUpdated.Email
	}
	if userUpdated.Url != "" {
		changes++
		mainUser.Url = userUpdated.Url
	}

	if changes == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("UserUpdate: ", mainUser)

	if err := repository.Data.UserUpdate(mainUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Start() {
	topRouter := mux.NewRouter()

	topRouter.HandleFunc("/login", LoginHandler).Methods("POST", "OPTIONS")
	topRouter.HandleFunc("/signup", SignupHandler).Methods("POST", "OPTIONS")

	userRouter := topRouter.PathPrefix("/api").Subrouter()

	userRouter.Use(middleware.AuthMiddleware)
	userRouter.HandleFunc("/user", GetProfileHandler).Methods("GET", "OPTIONS")
	userRouter.HandleFunc("/user", UpdateProfile).Methods("PUT", "OPTIONS")
	userRouter.HandleFunc("/user", LogoutHandler).Methods("DELETE", "OPTIONS")

	topRouter.Use(mux.CORSMethodMiddleware(topRouter))

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
