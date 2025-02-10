package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// .gitignore has object id: be22bdeea290f39e10c63f4f199a156ba1d9728a.
func main() {
	rawContent := "bin/\nvendor/"
	rawContentLength := len([]byte(rawContent))

	rawContentHasher := sha1.New()
	rawContentHasher.Write([]byte(rawContent))

	// prints c3f857d70bd30e3f504f1c6b89b5697930923f09
	fmt.Println("raw", hex.EncodeToString(rawContentHasher.Sum(nil)))

	content := fmt.Sprintf("blob %d\x00%s", rawContentLength, rawContent)

	hasher := sha1.New()
	hasher.Write([]byte(content))

	// prints be22bdeea290f39e10c63f4f199a156ba1d9728a
	fmt.Println("raw", hex.EncodeToString(hasher.Sum(nil)))
}
