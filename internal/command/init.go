package command

import (
	"errors"
	"fmt"
	"io"

	"github.com/LukasJenicek/ggit/internal/repository"
)

var ErrRepositoryAlreadyInitialized = errors.New("repository already initialized")

type InitCommand struct {
	repository *repository.Repository
}

func NewInitCommand(repo *repository.Repository) (*InitCommand, error) {
	if repo == nil {
		return nil, errors.New("repository is nil")
	}

	return &InitCommand{
		repository: repo,
	}, nil
}

func (i *InitCommand) Run() ([]byte, error) {
	if i.repository.Initialized {
		return nil, ErrRepositoryAlreadyInitialized
	}

	if err := i.repository.Init(); err != nil {
		return nil, fmt.Errorf("init repository: %w", err)
	}

	msg := "Initialized empty Git repository in " + i.repository.GitPath

	return []byte(msg), nil
}

func (i *InitCommand) Output(msg []byte, err error, stdout io.Writer) (int, error) {
	if err != nil {
		if errors.Is(err, ErrRepositoryAlreadyInitialized) {
			fmt.Fprintf(stdout, "Reinitialized existing Git repository in %s", i.repository.GitPath)

			return 0, nil
		}

		return 1, fmt.Errorf("run init: %w", err)
	}

	fmt.Fprintf(stdout, string(msg))

	return 0, nil
}
