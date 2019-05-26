package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/nsmith5/vgraas/pkg/middleware"
	"github.com/nsmith5/vgraas/pkg/vgraas"
)

func main() {
	var (
		addr = flag.String("api", ":8080", "API listen address")
	)
	flag.Parse()

	var api http.Handler
	{
		repo := vgraas.NewRAMRepo()
		api = vgraas.NewAPI(repo)

		// Limit request size to 500 KiB
		api = middleware.LimitBody(api, 1<<19)

		// Rate limit requests to 5Hz per remote address with bursts of 2
		api = middleware.RateLimit(api, 5, 2, middleware.XForwardedFor)
	}

	log.Println(http.ListenAndServe(*addr, api))
}
