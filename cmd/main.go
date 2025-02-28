package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/LukasJenicek/ggit/internal/clock"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/repository"
	"github.com/LukasJenicek/ggit/internal/workspace"
)

func main() {
	if len(os.Args) < 2 {
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
			if errors.Is(err, repository.ErrNoFilesToCommit) {
				fmt.Println(err)
				os.Exit(1)
			}

			log.Fatalf("commit: %v", err)
		}

		fmt.Printf("[%s] Successfully committed changes\n", cID)
	case "add":
		if len(os.Args) < 3 {
			fmt.Printf("usage: %s add <pattern>\n", os.Args[0])
			fmt.Println("Examples:")
			fmt.Println("  Add single file: ggit add file.txt")
			fmt.Println("  Add using glob pattern: ggit add \"*.go\"")

			os.Exit(0)
		}

		err := repo.Add(os.Args[2:])
		if err != nil {
			var cErr *workspace.ErrPathNotMatched
			if errors.As(err, &cErr) {
				fmt.Println(cErr)
				os.Exit(128)
			}

			log.Fatalf("add: %v", err)
		}
	default:
		log.Fatalf("unknown command: %q", cmd)
	}
}
