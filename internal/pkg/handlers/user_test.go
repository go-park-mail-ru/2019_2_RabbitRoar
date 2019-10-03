package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/entity"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/repository"
	"github.com/gorilla/context"
)

func TestUpdateProfile(t *testing.T) {
	expectedUser := entity.User{1, "Kekos", "Egos@mail.ru", 0, "", "password"}

	repository.Data.UserCreate("Kekos", "password", "Kek@mail.ru")
	user, err := repository.Data.UserGetByName("Kekos")
	if err != nil {
		t.Error("user was not created")
	}

	UUID, err := repository.Data.SessionCreate(user)
	if err != nil {
		t.Error("session was not created")

	}

	cookie := http.Cookie{}
	cookie.Name = "sessionID"
	cookie.Value = UUID.String()

	request_body := bytes.NewReader([]byte(`{"Email": "Egos@mail.ru"}`))

	r := httptest.NewRequest("PUT", "/user", request_body)
	r.AddCookie(&cookie)
	context.Set(r, "user", user)
	w := httptest.NewRecorder()

	t.Parallel()

	UpdateProfile(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	myUser, err := repository.Data.UserGetByName("Kekos")
	if err != nil {
		t.Error("user was not created")
	}

	reflect.DeepEqual(myUser, expectedUser)
}
