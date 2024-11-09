package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <cmd>\n", os.Args[0])
	}

	cmd := os.Args[1]
	if cmd == "" {
		log.Fatalf("cmd is required")
	}

	switch cmd {
	case "init":
		if err := initGit(); err != nil {
			log.Fatalf("init: %v", err)
		}
		fmt.Println("Initialized empty Git repository")
	default:
		log.Fatalf("unknown command: %q", cmd)
	}
}

func initGit() error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current working directory: %v", err)
	}

	paths := []string{".git", ".git/objects", ".git/refs"}
	for _, path := range paths {
		err = os.Mkdir(filepath.Join(workingDirectory, path), os.ModePerm)
		if err != nil {
			return fmt.Errorf("create %s directory: %w", path, err)
		}
	}

	return nil
}
