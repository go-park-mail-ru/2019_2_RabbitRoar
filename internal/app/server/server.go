package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test\n"))
}

func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", MainHandler)
	log.Fatal(http.ListenAndServe(":8000", r))
}
