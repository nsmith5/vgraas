package middleware

import "net/http"

// LimitBody is used to limit the amount of bytes in an http request that
// are read. This prevents malicious actors from uploaded massive requests.
func LimitBody(next http.Handler, bytes int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, bytes)
		next.ServeHTTP(w, r)
	})
}
