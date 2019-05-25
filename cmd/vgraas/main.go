package main

import (
	"os"

	"github.com/nsmith5/vgraas/pkg/vgraas"
)

func main() {
	vgraas.Hello(os.Stdout)
}
