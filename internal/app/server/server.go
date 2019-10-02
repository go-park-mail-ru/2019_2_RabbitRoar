package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Проверка командой curl --data "password=1234&username=Anita" http://localhost:3000/user/login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	inputUsername := r.FormValue("username")
	inputPassword := r.FormValue("password")
	fmt.Fprintln(w, "you username: ", inputUsername)
	fmt.Fprintln(w, "you password: ", inputPassword)

}
func SignupHandler(w http.ResponseWriter, r *http.Request) {
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
}

func SetProfile(w http.ResponseWriter, r *http.Request) {
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
}

func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/user/login", LoginHandler).Methods("POST")
	r.HandleFunc("/user/signup", SignupHandler).Methods("POST")
	r.HandleFunc("/user", SetProfile).Methods("GET")
	r.HandleFunc("/user", UpdateProfile).Methods("PUT")
	r.HandleFunc("/user", LogoutHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", r))
}
