package usecase

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/utils"
	"github.com/microcosm-cc/bluemonday"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var log = logging.MustGetLogger("user_handler")

type userUseCase struct {
	repository user.Repository
	sanitizer *bluemonday.Policy
}

func NewUserUseCase(userRepo user.Repository) user.UseCase {
	return &userUseCase{
		repository: userRepo,
		sanitizer: bluemonday.UGCPolicy(),
	}
}

func (uc *userUseCase) Sanitize(u models.User) models.User {
	u.Username = uc.sanitizer.Sanitize(u.Username)
	u.Email = uc.sanitizer.Sanitize(u.Email)
	u.Password = ""
	return u
}

func (uc *userUseCase) Create(u models.User) (*models.User, error) {
	if err := uc.prepare(&u); err != nil {
		return nil, err
	}

	return uc.repository.Create(u)
}

func (uc *userUseCase) prepare(u *models.User) error {
	ok, err := govalidator.ValidateStruct(u)
	if !ok {
		return err
	}

	if err := uc.preparePassword(u); err != nil {
		return err
	}

	if err := uc.prepareUsername(u); err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) preparePassword(u *models.User) error {
	u.Password = auth.HashPassword(u.Password)
	return nil
}

func (uc *userUseCase) prepareUsername(u *models.User) error {
	return nil
}

func (uc *userUseCase) Update(u, uUpdate models.User) (*models.User, error) {
	if uUpdate.Password != "" {
		u.Password = uUpdate.Password
	}

	if uUpdate.Username != "" {
		u.Username = uUpdate.Username
	}

	if err := uc.prepare(&u); err != nil {
		return nil, err
	}

	return &u, uc.repository.Update(u)
}

func (uc *userUseCase) UpdateAvatar(u models.User, file *multipart.FileHeader) (*models.User, error) {
	src, err := file.Open()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer src.Close()

	ok, ext := checkFileContentType(src)
	if !ok {
		log.Errorf("error invalid contentType")
		return nil, errors.New("error Invalid ContentType")
	}

	filename := strconv.Itoa(u.ID) + "." + ext

	// Move to config
	filePath := filepath.Join(
		"data",
		"uploads",
		"avatar",
		filename,
	)

	dst, err := os.Create(filePath)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Error(err)
		return nil, err
	}

	url := viper.GetString("server.static.avatar_prefix") + filename
	u.AvatarUrl = url

	return &u, uc.repository.Update(u)
}

func checkFileContentType(file multipart.File) (bool, string) {
	contentType, err := getFileContentType(file)
	if err != nil {
		return false, ""
	}
	allowedContentType := []string{
		"image/png",
		"image/jpeg",
	}
	if utils.SliceContains(allowedContentType, contentType) {
		return true, strings.Split(contentType, "/")[1]
	}
	return false, ""
}

func getFileContentType(file multipart.File) (string, error){
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func (uc *userUseCase) GetByID(id int) (*models.User, error) {
	return uc.repository.GetByID(id)
}

func (uc *userUseCase) GetByName(name string) (*models.User, error) {
	return uc.repository.GetByName(name)
}

func (uc *userUseCase) GetBySessionID(UUID string) (*models.User, error) {
	return uc.repository.GetBySessionID(UUID)
}

func (uc *userUseCase) IsPasswordCorrect(u models.User) (*models.User, bool) {
	correctUser, err := uc.repository.GetByName(u.Username)
	if err != nil {
		return nil, false
	}

	if !auth.CheckPassword(u.Password, correctUser.Password) {
		return nil, false
	}

	return correctUser, true
}
