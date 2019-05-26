package main

import (
	"flag"
	"log"
	"net/http"

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
	}

	log.Println(http.ListenAndServe(*addr, api))
}
