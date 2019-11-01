package auth

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) string {
	passwordBytes := []byte(password)

	salt := make([]byte, 8)
	_, _ = rand.Read(salt)

	hashedPasswordBytes := hashPassword(passwordBytes, salt)
	hashedPassword := string(hashedPasswordBytes)

	return hashedPassword
}

func hashPassword(password, salt []byte) []byte {
	hashedPassword := argon2.IDKey(
		password,
		salt,
		1,
		64*1024,
		4,
		32)
	return append(salt, hashedPassword...)
}

func CheckPassword(userPassword, passwordHashed string) bool {
	passwordBytesHashed := []byte(passwordHashed)
	salt := passwordBytesHashed[0:8]

	userPasswordBytes := []byte(userPassword)
	userPasswordBytesHashed := hashPassword(userPasswordBytes, salt)

	return bytes.Equal(userPasswordBytesHashed, passwordBytesHashed)
}
