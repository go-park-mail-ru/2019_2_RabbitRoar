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
		sessionIDCookie, errCookie := r.Cookie("SessionID")

		fmt.Println("Middleware: Got request: ", r)

		fmt.Println("Middleware: got cookie: ", sessionIDCookie)

		if errCookie != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		fmt.Println("Middleware: ", errCookie)

		sessionID, errParseUUID := uuid.Parse(sessionIDCookie.Value)

		if errParseUUID != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		fmt.Println("Middleware: ", errParseUUID)

		user, err := repository.Data.SessionGetUser(sessionID)

		fmt.Println("Middleware: ", err)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		context.Set(r, "user", user)

		log.Printf("Authenticated user %s\n", user)
		next.ServeHTTP(w, r)
	})
}
