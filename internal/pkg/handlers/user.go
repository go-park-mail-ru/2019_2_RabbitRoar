package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/entity"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/repository"
	"github.com/gorilla/context"
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
