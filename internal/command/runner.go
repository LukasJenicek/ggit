package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"io"
	"os"
	"path/filepath"

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
	cmd, err := NewCommitCmd(r.repository)
	if err != nil {
		return 1, fmt.Errorf("init commit cmd: %w", err)
	}

	out, err := cmd.Run()
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
		help := `Usage: ggit add <pattern>
Examples:
	Add single file: ggit add file.txt
	Add using glob pattern: ggit add *.go
`
		fmt.Fprintf(output, help)

		return 0, nil
	}

	cmd, err := NewAddCommand(args, r.repository)
	if err != nil {
		return 1, fmt.Errorf("init add command: %w", err)
	}

	out, err := cmd.Run()
	if err != nil {
		var cErr *workspace.ErrPathNotMatched
		if errors.As(err, &cErr) {
			fmt.Fprintf(output, "%s", cErr.Error())
			return 128, nil
		}

		if errors.Is(err, os.ErrPermission) {
			fmt.Fprintf(output, "error: permission denied\n")
			fmt.Fprintf(output, "fatal: adding files failed")
			return 128, nil
		}

		if errors.Is(err, filesystem.ErrLockAcquired) {
			help := `Another git process seems to be running in this repository, e.g.
an editor opened by 'git commit'. Please make sure all processes
are terminated then try again. If it still fails, a git process
may have crashed in this repository earlier:
remove the file manually to continue.`
			fmt.Fprintf(
				output,
				"fatal: Unable to create %s: File exists\n",
				filepath.Join(r.repository.GitPath, "index.lock"),
			)
			fmt.Fprintf(output, help)
			return 128, nil
		}

		return 1, fmt.Errorf("add files: %w", err)
	}

	fmt.Fprintf(output, "%s", string(out))

	return 0, nil
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
