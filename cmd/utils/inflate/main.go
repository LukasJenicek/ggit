package main

import (
	"compress/zlib"
	"errors"
	"io"
	"log"
	"os"
)

// flate compression algorithm (used in gzip, zlib, and raw DEFLATE)
// gzip starts 0x1F 0x8B
//
// zlib headers:
// 0x78 0x01  (Fastest)
// 0x78 0x9C  (Default)
// 0x78 0xDA  (Best compression).
//
// raw deflate does not have any headers
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
