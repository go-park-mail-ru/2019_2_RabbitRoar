package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/entity"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/repository"
	"github.com/gorilla/context"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("handlers")

func GetProfileHandler(ctx *echo.Context) {

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
	if userUpdated.AvatarUrl != "" {
		changes++
		mainUser.AvatarUrl = userUpdated.AvatarUrl
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

var allowedContentType = [...]string{"image/png", "image/jpeg"}

func UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(entity.User)

	var err error

	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		logger.Error("Error parsing multipart form: ", err)
		return
	}

	tmpFile, _, err := r.FormFile("uploadfile")
	if err != nil {
		logger.Error("Error parsing multipart form: ", err)
		return
	}
	defer tmpFile.Close()

	// 512 count of bytes needed at most to detect Content-Type
	detectContentTypeBuffer := make([]byte, 512)
	_, err = tmpFile.Read(detectContentTypeBuffer)

	contentType := http.DetectContentType(detectContentTypeBuffer)

	if !utils.Contains(allowedContentType, contentType) {
		logger.Warning("Image type not allowed.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = tmpFile.Seek(0, 0)

	// Move to config
	file, err := os.Create("data/uploads/" + user.Username)
	if err != nil {
		logger.Error("Error creating file: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, tmpFile)
	if err != nil {
		logger.Error("Error copy file: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.AvatarUrl = "data/uploads/" + user.Username

	w.WriteHeader(http.StatusOK)
	return
}
