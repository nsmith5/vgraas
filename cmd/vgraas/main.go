package main

import (
	"io"
	"log"
	"net/http"

	"github.com/nsmith5/vgraas/pkg/vgraas"
)

func HelloHTTP(w http.ResponseWriter, r *http.Request) {
	err := vgraas.Hello(w)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "Failed to write hello :(")
		return
	}
	return
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloHTTP)
	log.Println(http.ListenAndServe(":8080", mux))
}
