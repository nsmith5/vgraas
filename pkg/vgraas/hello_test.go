package vgraas

import (
	"strings"
	"testing"
)

func TestHello(t *testing.T) {
	var b strings.Builder
	err := Hello(&b)
	if err != nil {
		t.Error(err)
	}

	if n := strings.Compare("Hello, World!\n", b.String()); n != 0 {
		t.Error("Wrote the wrong string")
	}
}
