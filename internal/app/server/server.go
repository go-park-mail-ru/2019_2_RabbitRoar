package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mux/context"
	"net/http"

	"../../../cmd/entity"
	"../../../cmd/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//  curl -XPOST -H "Content-type: application/json" -d '{"username":"anita", "password":"1234"}' 'http://localhost:3000/user/login'
func LoginHandler(w http.ResponseWriter, r *http.Request) {
}

// curl -XPOST -H "Content-type: application/json" -d '{"username":"anita", "password":"1234", "email":"anit@mail.com"}' 'http://localhost:3000/user/signup'
func SignupHandler(w http.ResponseWriter, r *http.Request) {
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

	if cookie, err := r.Cookie("sessionID"); err != nil {
		if UUID, err := uuid.Parse(cookie.Value); err != nil {
		} else {
			repository.Data.SessionDestroy(UUID)
		}
	}
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := context.Get("user")

	userJSON, err := json.Marshal(user)

	if err != nil {
		panic("error marshaling user")
	}

	fmt.Fprintln(user, userJSON)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
}

func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/user/login", LoginHandler).Methods("POST")
	r.HandleFunc("/user/signup", SignupHandler).Methods("POST")
	r.HandleFunc("/user", GetProfileHandler).Methods("GET")
	r.HandleFunc("/user", UpdateProfile).Methods("PUT")
	r.HandleFunc("/user", LogoutHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", r))
}
