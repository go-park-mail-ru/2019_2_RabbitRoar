package csrf

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtToken struct {
	Secret []byte
}

type JwtCsrfClaims struct {
	SessionID string
	jwt.StandardClaims
}

func (tk *JwtToken) Create(s string) (string, error) {
	data := JwtCsrfClaims{
		SessionID: s,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour*24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(tk.Secret)
}

func (tk *JwtToken) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return tk.Secret, nil
}

func (tk *JwtToken) Check(s string, inputToken string) (bool, error) {
	payload := &JwtCsrfClaims{}
	_, err := jwt.ParseWithClaims(inputToken, payload, tk.parseSecretGetter)
	if err != nil {
		return false, fmt.Errorf("cant parse jwt token: %v", err)
	}
	if payload.Valid() != nil {
		return false, fmt.Errorf("invalid jwt token: %v", err)
	}
	return payload.SessionID == s, nil
}
