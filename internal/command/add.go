package command

import (
	"fmt"
	"github.com/LukasJenicek/ggit/internal/repository"
)

type AddCommand struct {
	paths      []string
	repository *repository.Repository
}

func NewAddCommand(paths []string, repository *repository.Repository) (*AddCommand, error) {
	if repository == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("paths is empty")
	}

	allEmpty := true
	for _, path := range paths {
		if path != "" {
			allEmpty = false
		}
	}

	if allEmpty {
		return nil, fmt.Errorf("paths is empty")
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
