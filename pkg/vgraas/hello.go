package vgraas

import (
	"io"
)

func Hello(w io.Writer) (err error) {
	_, err = io.WriteString(w, "Hello, World!\n")
	return err
}
