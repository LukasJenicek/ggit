package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	path := "libs/internal/haha.txt"

	fmt.Println(filepath.Dir(path))
}
