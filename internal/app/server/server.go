package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/entity"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/middleware"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//  curl -XPOST -H "Content-type: application/json" -d '{"username":"anita", "password":"1234"}' 'http://localhost:3000/user/login'
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	u := &entity.User{}
	json.Unmarshal(body, u)
	if u.Username == "" || u.Password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err = repository.Data.UserGetByName(u.Username)
	if err == repository.ErrNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := repository.Data.UserGetByName(u.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if user.Password != u.Password { //SORRY!!!
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid, _ := repository.Data.SessionCreate(user)
	SessionID := uuid.String()

	http.SetCookie(w, &http.Cookie{
		Name:     "SessionID",
		Value:    SessionID,
		Path:     "localhost:3000",
		HttpOnly: true,
		// Secure: true,
	})

	w.WriteHeader(http.StatusOK)
	return
}

// curl -XPOST -H "Content-type: application/json" -d '{"username":"anita", "password":"1234", "email":"anit@mail.com"}' 'http://localhost:3000/user/signup'
func SignupHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	u := &entity.User{}
	json.Unmarshal(body, u)
	if u.Username == "" || u.Password == "" || u.Email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err = repository.Data.UserCreate(u.Username, u.Password, u.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		w.WriteHeader(http.StatusOK)
		r.Body.Close()
	}()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	if cookie, err := r.Cookie("SessionID"); err != nil {
		if UUID, err := uuid.Parse(cookie.Value); err != nil {
		} else {
			repository.Data.SessionDestroy(UUID)
		}
	}
}

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
	topRouter.HandleFunc("/login", LoginHandler).Methods("POST")
	topRouter.HandleFunc("/signup", SignupHandler).Methods("POST")

	userRouter := topRouter.PathPrefix("/api").Subrouter()
	userRouter.Use(middleware.AuthMiddleware)
	userRouter.HandleFunc("/user", GetProfileHandler).Methods("GET")
	userRouter.HandleFunc("/user", UpdateProfile).Methods("PUT")
	userRouter.HandleFunc("/user", LogoutHandler).Methods("DELETE")

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
