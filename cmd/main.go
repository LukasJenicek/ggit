package main

import (
	"fmt"
	"log"
	"os"

	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/repository"
	"github.com/LukasJenicek/ggit/internal/workspace"
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
		w := workspace.New()

		files, err := w.ListFiles(workingDirectory)
		if err != nil {
			log.Fatalf("list files: %v", err)
		}
		fmt.Println(files)
	default:
		log.Fatalf("unknown command: %q", cmd)
	}
}
