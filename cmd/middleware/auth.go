package middleware

import (
	"log"
	"mux/context"
	"net/http"

	"../repository"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Cookie("sessionID")

		user, err := repository.Data.SessionGetUser(sessionID)

		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}

		context.Set(r, "user", user)

		log.Printf("Authenticated user %s\n", user)
		next.ServeHTTP(w, r)
	})
}
