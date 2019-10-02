package middleware

import (
	"fmt"
	"log"
	"net/http"

	"../repository"
	"github.com/google/uuid"
	"github.com/gorilla/context"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionIDCookie, errCookie := r.Cookie("sessionID")

		if errCookie != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		fmt.Println(errCookie)

		sessionID, errParseUUID := uuid.Parse(sessionIDCookie.Value)

		if errParseUUID != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		fmt.Println(errParseUUID)

		user, err := repository.Data.SessionGetUser(sessionID)

		fmt.Println(err)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		context.Set(r, "user", user)

		log.Printf("Authenticated user %s\n", user)
		next.ServeHTTP(w, r)
	})
}
