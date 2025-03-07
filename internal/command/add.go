package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/repository"
	"github.com/LukasJenicek/ggit/internal/workspace"
)

var LockAcquiredMsg = `Another git process seems to be running in this repository, e.g.
an editor opened by 'git commit'. Please make sure all processes
are terminated then try again. If it still fails, a git process
may have crashed in this repository earlier:
remove the file manually to continue.
`

type AddCommand struct {
	paths      []string
	repository *repository.Repository
}

func NewAddCommand(paths []string, repository *repository.Repository) (*AddCommand, error) {
	if repository == nil {
		return nil, errors.New("repository is nil")
	}

	if len(paths) == 0 {
		return nil, errors.New("paths is empty")
	}

	allEmpty := true

	for _, path := range paths {
		if path != "" {
			allEmpty = false
		}
	}

	if allEmpty {
		return nil, errors.New("paths is empty")
	}

	return &AddCommand{
		paths:      paths,
		repository: repository,
	}, nil
}

func (a *AddCommand) Run() ([]byte, error) {
	if err := a.repository.Add(a.paths); err != nil {
		return nil, err
	}

	return []byte(""), nil
}

func (a *AddCommand) Output(msg []byte, err error, stdout io.Writer) (int, error) {
	if err != nil {
		var cErr *workspace.ErrPathNotMatched
		if errors.As(err, &cErr) {
			fmt.Fprintf(stdout, "%s", cErr.Error())

			return 128, nil
		}

		if errors.Is(err, os.ErrPermission) {
			fmt.Fprintf(stdout, "error: permission denied\n")
			fmt.Fprintf(stdout, "fatal: adding files failed")

			return 128, nil
		}

		if errors.Is(err, filesystem.ErrLockAcquired) {
			fmt.Fprintf(
				stdout,
				"fatal: Unable to create %s: File exists\n",
				filepath.Join(a.repository.GitPath, "index.lock"),
			)
			fmt.Fprintf(stdout, "%s", LockAcquiredMsg)

			return 128, nil
		}

		return 1, fmt.Errorf("add files: %w", err)
	}

	fmt.Fprint(stdout, string(msg))

	return 0, nil
}
