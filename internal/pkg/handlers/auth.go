package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/entity"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/repository"
	"github.com/google/uuid"
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
