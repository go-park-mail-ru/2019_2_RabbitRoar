package server

import (
	"../../../cmd/entity"
	"../../../cmd/repository"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// curl --data "password=1234&username=Anita" http://localhost:3000/user/login
//
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	u := entity.User{}
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}
	//b, err := ioutil.ReadAll(r.Body)
	//fmt.Println(b)
	//defer r.Body.Close()
	//
	//if err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
	//
	//u := entity.User{}
	//err = json.Unmarshal(b, &u)
	//if err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
	//fmt.Fprintln(w, u)

	//output, err := json.Marshal(u)
	//if err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
	//w.Header().Set("content-type", "application/json")
	//w.Write(output)
}

// curl --data "password=1234&username=Anita&email=heh@mail.ru" http://localhost:3000/user/signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	inputUsername := r.FormValue("username")
	inputPassword := r.FormValue("password")
	inputEmail := r.FormValue("email")
	createdUser, err := repository.Data.UserCreate(inputUsername, inputPassword, inputEmail)
	if err != nil {
		fmt.Fprintln(w, "creation error")
	}
	fmt.Fprintln(w, createdUser)
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
