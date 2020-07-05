package middlewares

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/vishal979/auth/server/auth"
	"github.com/vishal979/auth/server/responses"
)

// SetMiddlewareJSON set output to json
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// SetMiddlewareAuthentication sets middleware to validate token
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			log.Error("invalid token access")
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
