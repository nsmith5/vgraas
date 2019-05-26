package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// The interceptingWriter is stolen from the very clever Peter Bourgon.
// The code is Apache 2 licensed and can be found here:
// https://github.com/peterbourgon/breakfast-solutions/blob/master/misc.go

// interceptingWriter records the count of bytes writen to a Writer and the
// status code. Handy for logging what happens to an http request without
// injecting the logging code into the request handling logic.
type interceptingWriter struct {
	count int
	code  int
	http.ResponseWriter
}

func (iw *interceptingWriter) WriteHeader(code int) {
	iw.code = code
	iw.ResponseWriter.WriteHeader(code)
}

func (iw *interceptingWriter) Write(p []byte) (int, error) {
	iw.count += len(p)
	return iw.ResponseWriter.Write(p)
}

// Logging is a middleware that adds structured logging (JSON) to an
// http.Handler.
//
// Time, method, path, status, response size and duration are recorded
// in newline delimited JSON documents in the supplied io.Writer.
func Logging(next http.Handler, out io.Writer) http.Handler {
	enc := json.NewEncoder(out)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		iw := &interceptingWriter{0, http.StatusOK, w}
		now := time.Now()

		next.ServeHTTP(iw, r)

		enc.Encode(map[string]interface{}{
			"time":     now,
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   iw.code,
			"respSize": iw.count,
			"duration": time.Since(now).Seconds(),
		})
	})
}
