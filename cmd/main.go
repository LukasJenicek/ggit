package main

import (
	"fmt"
	"log"
	"os"

	"github.com/LukasJenicek/ggit/internal/clock"
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

	repo, err := repository.New(filesystem.New(), &clock.RealClock{}, workingDirectory)
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

		fmt.Println("Initialized empty Git repository in", repo.GitPath)
	case "commit":
		cID, err := repo.Commit()
		if err != nil {
			log.Fatalf("commit: %v", err)
		}
		fmt.Printf("[%s] Successfully committed changes\n", cID)
	default:
		log.Fatalf("unknown command: %q", cmd)
	}
}
