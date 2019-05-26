package middleware

import "net/http"

// ContentType is used to set the Content-Type header for http handlers.
//
// This can be useful to set the Content-Type for an entire mux.Router, for example.
func ContentType(next http.Handler, header string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", header)
		next.ServeHTTP(w, r)
	})
}
