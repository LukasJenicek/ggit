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

// Run
// It shows a summary of the differences between the tree of the HEAD commit, the entries in the index, and
// the contents of the workspace, as well as listing conflicting files during a merge
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
	untrackedFiles := []string{}
	trackedFiles := []string{}

	for _, f := range files {
		if _, ok := entries[f]; ok {
			trackedFiles = append(trackedFiles, f)
			continue
		}

		untrackedFiles = append(untrackedFiles, f)
		buf.WriteString(fmt.Sprintf("?? %s\n", f))
	}

	return buf.Bytes(), nil
}

func (s *StatusCommand) Output(msg []byte, err error, stdout io.Writer) (int, error) {
	if err != nil {
		return 1, fmt.Errorf("output: %w", err)
	}

	fmt.Fprint(stdout, string(msg))

	return 0, nil
}
