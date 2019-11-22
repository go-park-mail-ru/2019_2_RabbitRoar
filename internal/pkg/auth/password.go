package auth

import (
	"bytes"
	"crypto/rand"
	"github.com/prometheus/common/log"
	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) string {
	salt := make([]byte, 8)
	_, _ = rand.Read(salt)

	hashedPassword := hashPassword([]byte(password), salt)

	return string(hashedPassword)
}

func hashPassword(password, salt []byte) []byte {
	hashedPassword := argon2.IDKey(
		password,
		salt,
		1,
		64*1024,
		4,
		32,
	)
	return append(hashedPassword, salt...)
}

func CheckPassword(userPassword, passwordHashed string) bool {
	passwordHashedBytes := []byte(passwordHashed)
	salt := passwordHashedBytes[len(passwordHashed)-8:]

	userPasswordBytes := []byte(userPassword)
	userPasswordBytesHashed := hashPassword(userPasswordBytes, salt)

	log.Info("Invalid password")
	return bytes.Equal(userPasswordBytesHashed, passwordHashedBytes)
}
