package command

import (
	"bytes"
	"fmt"
	"io"

	"github.com/LukasJenicek/ggit/internal/repository"
)

type StatusCommand struct {
	repo *repository.Repository
}

func NewStatusCommand(repo *repository.Repository) (*StatusCommand, error) {
	return &StatusCommand{repo}, nil
}

func (s *StatusCommand) Run() ([]byte, error) {
	files, err := s.repo.Workspace.ListFiles()
	if err != nil {
		return nil, fmt.Errorf("list files: %w", err)
	}

	entries, _, err := s.repo.Index.LoadEntries()
	if err != nil {
		return nil, fmt.Errorf("load index entries: %w", err)
	}

	buf := bytes.NewBuffer(nil)

	for _, f := range files {
		if _, ok := entries[f]; ok {
			continue
		}

		buf.WriteString(fmt.Sprintf("?? %s\n", f))
	}

	return buf.Bytes(), nil
}

func (s *StatusCommand) Output(msg []byte, err error, stdout io.Writer) (int, error) {
	if err != nil {
		return 1, fmt.Errorf("output: %w", err)
	}

	fmt.Fprint(stdout, "%s", string(msg))

	return 0, nil
}
