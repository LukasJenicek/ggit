package command

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/LukasJenicek/ggit/internal/repository"
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
	case "add":
		return r.addCmd(args, output)
	case "commit":
		return r.commitCmd(output)
	case "init":
		return r.initCmd(output)
	}

	return 1, fmt.Errorf("ggit: %q is not a ggit command. See 'ggit --help'", cmd)
}

func (r *Runner) commitCmd(output io.Writer) (int, error) {
	cmd, err := NewCommitCmd(r.repository)
	if err != nil {
		return 1, fmt.Errorf("init commit cmd: %w", err)
	}

	out, err := cmd.Run()

	return cmd.Output(out, err, output)
}

func (r *Runner) addCmd(args []string, output io.Writer) (int, error) {
	if len(args) == 0 {
		help := `Usage: ggit add <pattern>
Examples:
	Add single file: ggit add file.txt
	Add using glob pattern: ggit add *.go
`
		fmt.Fprint(output, help)

		return 0, nil
	}

	cmd, err := NewAddCommand(args, r.repository)
	if err != nil {
		return 1, fmt.Errorf("init add command: %w", err)
	}

	out, err := cmd.Run()

	return cmd.Output(out, err, output)
}

func (r *Runner) initCmd(output io.Writer) (int, error) {
	cmd, err := NewInitCommand(r.repository)
	if err != nil {
		return 1, fmt.Errorf("init command: %w", err)
	}

	out, err := cmd.Run()
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
