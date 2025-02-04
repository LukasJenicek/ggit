package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/repository"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <cmd>\n", os.Args[0])
	}

	cmd := os.Args[1]
	if cmd == "" {
		log.Fatalf("cmd is required")
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current working directory: %v", err)
	}

	repo, err := repository.New(filesystem.New(), workingDirectory)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}

	switch cmd {
	case "init":
		if repo.Initialized {
			log.Fatalf("ggit already initialized")
		}

		if err = repo.Init(); err != nil {
			log.Fatalf("ggit init: %v", err)
		}
	case "commit":
		files, err := listFiles(workingDirectory)
		if err != nil {
			log.Fatalf("list files: %v", err)
		}
		fmt.Println(files)
	default:
		log.Fatalf("unknown command: %q", cmd)
	}
}

func listFiles(workingDir string) ([]string, error) {
	// TODO: load more ignored files from config
	ignore := []string{".", "..", ".git", workingDir}

	files := []string{}
	err := filepath.WalkDir(workingDir, func(path string, d fs.DirEntry, err error) error {
		if slices.Contains(ignore, path) {
			return nil
		}

		if slices.Contains(ignore, d.Name()) {
			return filepath.SkipDir
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk recursively dir %q: %v", workingDir, err)
	}

	return files, nil
}
