package command

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/LukasJenicek/ggit/internal/repository"
	"github.com/LukasJenicek/ggit/internal/workspace"
)

type Runner struct {
	repository *repository.Repository
}

func NewRunner(repo *repository.Repository) *Runner {
	return &Runner{
		repository: repo,
	}
}

// RunCmd (osExit, err).
func (r *Runner) RunCmd(ctx context.Context, cmd string, args []string, output io.Writer) (int, error) {
	switch cmd {
	case "init":
		return r.initCmd(output)
	case "commit":
		return r.commitCmd(output)
	case "add":
		return r.addCmd(args, output)
	}

	return 127, fmt.Errorf("unknown command: %s", cmd)
}

func (r *Runner) commitCmd(output io.Writer) (int, error) {
	commitCmd, err := NewCommitCmd(r.repository)
	if err != nil {
		return 1, fmt.Errorf("init commit cmd: %w", err)
	}

	out, err := commitCmd.Run()
	if err != nil {
		if errors.Is(err, repository.ErrNoFilesToCommit) {
			fmt.Fprintf(output, "%s", err.Error())

			return 1, nil
		}

		return 1, fmt.Errorf("commit cmd: %w", err)
	}

	fmt.Fprintf(output, "%s", string(out))

	return 0, nil
}

func (r *Runner) addCmd(args []string, output io.Writer) (int, error) {
	if len(args) == 0 {
		fmt.Fprintf(output, "usage: %s add <pattern>\n", "add")
		fmt.Fprintf(output, "Examples:\n")
		fmt.Fprintf(output, "  Add single file: ggit add file.txt\n")
		fmt.Fprintf(output, "  Add using glob pattern: ggit add \"*.go\"\n")

		return 0, nil
	}

	addCmd, err := NewAddCommand(args, r.repository)
	if err != nil {
		return 1, fmt.Errorf("init add command: %w", err)
	}

	out, err := addCmd.Run()
	if err != nil {
		var cErr *workspace.ErrPathNotMatched
		if errors.As(err, &cErr) {
			fmt.Fprintf(output, "%s", err.Error())

			return 128, nil
		}

		return 1, fmt.Errorf("add files: %w", err)
	}

	fmt.Fprintf(output, "%s", string(out))

	return 0, nil
}

func (r *Runner) initCmd(output io.Writer) (int, error) {
	initCommand, err := NewInitCommand(r.repository)
	if err != nil {
		return 1, fmt.Errorf("init command: %w", err)
	}

	out, err := initCommand.Run()
	if err != nil {
		if errors.Is(err, ErrRepositoryAlreadyInitialized) {
			fmt.Fprintf(output, "Reinitialized existing Git repository in %s", r.repository.GitPath)

			return 0, nil
		}

		return 1, fmt.Errorf("run init: %w", err)
	}

	fmt.Fprintf(output, "%s", string(out))

	return 0, nil
}
