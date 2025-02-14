package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	path := "libs/haha.txt"

	fmt.Println(filepath.Dir(path))
	fmt.Println(strings.Split(path, "/"))
}
