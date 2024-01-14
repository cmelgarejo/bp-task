package middleware

import (
	"net/http"
	"os"
)

func BasicAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !checkUsernameAndPassword(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func checkUsernameAndPassword(username, password string) bool {
	return username == os.Getenv("BASIC_AUTH_USER") && password == os.Getenv("BASIC_AUTH_PASS")
}
