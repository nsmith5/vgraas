// Rate limiting code from https://www.alexedwards.net/blog/how-to-rate-limit-http-requests
// Licensed as MIT
package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Options for find the address to use for IP address based
// rate limiting.
const (
	// XForwardedFor uses the 'X-Forwarded-For' header to get the
	// ip address. This is useful behind reverse proxies.
	XForwardedFor int = iota

	// XRealIP uses the 'X-Real-IP' header to get ip addresses. Again
	// useful if behind a reverse proxy.
	XRealIP

	// RemoteAddr uses the requests raw remote address. Useful
	// if the API is exposes to direct connections.
	RemoteAddr
)

// RateLimit is a middleware that implements rate limiting based on
// IP address.
//
// It uses token-bucket algorithm with 'r' as rate and 'b' as burst.
// 'method' is used to specify the method for collecting the IP address.
func RateLimit(next http.Handler, r, b, method int) http.Handler {
	var visitors = make(map[string]*visitor)
	var mtx sync.Mutex

	// Periodically clean out stale entries forever
	// TODO: This is mildly irresponsible. We should have
	// clean shutdown on this thing.
	go func() {
		for {
			time.Sleep(time.Minute)
			mtx.Lock()
			for ip, v := range visitors {
				if time.Now().Sub(v.lastSeen) > 10*time.Minute {
					delete(visitors, ip)
				}
			}
			mtx.Unlock()
		}
	}()

	addVisitor := func(ip string) *rate.Limiter {
		limiter := rate.NewLimiter(rate.Limit(r), b)
		mtx.Lock()
		visitors[ip] = &visitor{limiter, time.Now()}
		mtx.Unlock()
		return limiter
	}

	getVisitor := func(ip string) *rate.Limiter {
		mtx.Lock()
		v, exists := visitors[ip]
		if !exists {
			mtx.Unlock()
			return addVisitor(ip)
		}

		v.lastSeen = time.Now()
		mtx.Unlock()
		return v.limiter
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var limiter *rate.Limiter
		switch method {
		case XForwardedFor:
			limiter = getVisitor(r.Header.Get("X-Forwarded-For"))
		case XRealIP:
			limiter = getVisitor(r.Header.Get("X-Real-IP"))
		case RemoteAddr:
			limiter = getVisitor(r.RemoteAddr)
		default:
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
