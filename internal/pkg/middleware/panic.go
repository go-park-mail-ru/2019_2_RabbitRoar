package middleware

import (
	"net/http"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("middleware_panic")

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			logger.Critical(err)
		}()
		next.ServeHTTP(w, r)
	})
}
