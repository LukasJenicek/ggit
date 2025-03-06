package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/LukasJenicek/ggit/internal/clock"
	"github.com/LukasJenicek/ggit/internal/command"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/repository"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <cmd>\n", os.Args[0])
	}

	cmd := os.Args[1]
	if cmd == "" {
		log.Fatalf("cmd is required")
	}

	ctx := context.Background()

	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current working directory: %v", err)
	}

	repo, err := repository.New(filesystem.New(), &clock.RealClock{}, workingDirectory)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}

	if cmd != "init" && !repo.Initialized {
		fmt.Println("fatal: not a git repository (or any of the parent directories): .git")
		os.Exit(128)
	}

	runner := command.NewRunner(repo)
	osExit, err := runner.RunCmd(ctx, cmd, os.Args[2:], os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(osExit)
	}

	os.Exit(osExit)
}
