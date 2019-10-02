package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"../../../cmd/repository"
	"../../../cmd/entity"
)

//  curl -XPOST -H "Content-type: application/json" -d '{"username": "anita", "password":"1234", "email":"anit@mail.com"}' 'http://localhost:3000/user/login'
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close() // важный пункт!
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	u := &entity.User{}
	json.Unmarshal(body, u)
	createdUser, err := repository.Data.UserCreate(u.Username, u.Password, u.Email)
	if err != nil {
		fmt.Fprintln(w, "creation error")
	}
	fmt.Fprintln(w, createdUser)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
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
