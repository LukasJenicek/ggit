package main

import (
	"compress/zlib"
	"errors"
	"io"
	"log"
	"os"
)

// raw deflate does not have any headers.
func main() {
	r, err := zlib.NewReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	for {
		if _, err = io.CopyN(os.Stdout, r, 1024); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			if errors.Is(err, io.EOF) {
				log.Fatal(err)
			}
		}
	}
}
