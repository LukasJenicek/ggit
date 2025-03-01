package command

import (
	"errors"
	"fmt"
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

	msg := fmt.Sprintf("Initialized empty Git repository in %s", i.repository.GitPath)

	return []byte(msg), nil
}
