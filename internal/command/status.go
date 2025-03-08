package command

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/LukasJenicek/ggit/internal/ds"
	"github.com/LukasJenicek/ggit/internal/repository"
)

// Shows difference between the tree of the HEAD commit, the entries in the index and workspace content.
type StatusCommand struct {
	repo *repository.Repository
}

func NewStatusCommand(repo *repository.Repository) (*StatusCommand, error) {
	return &StatusCommand{repo}, nil
}

func (s *StatusCommand) Run() ([]byte, error) {
	// Load entries into memory
	_, _, err := s.repo.Index.LoadEntries()
	if err != nil {
		return nil, fmt.Errorf("load index entries: %w", err)
	}

	untrackedFiles, err := s.scanWorkspace("")
	if err != nil {
		return nil, fmt.Errorf("scan workspace: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	for _, f := range untrackedFiles {
		buf.WriteString(fmt.Sprintf("?? %s\n", f))
	}

	return buf.Bytes(), nil
}

func (s *StatusCommand) scanWorkspace(prefix string) ([]string, error) {
	stats, err := s.repo.Workspace.ListDir(prefix)
	if err != nil {
		return nil, fmt.Errorf("list dir: %w", err)
	}

	untrackedFiles := ds.Set[string]{}

	for _, stat := range stats {
		path := stat.RelPath

		if s.repo.Index.Tracked(path) {
			if stat.FileInfo.IsDir() {
				if prefix != "" {
					path = filepath.Join(prefix, path)
				}

				files, err := s.scanWorkspace(path)
				if err != nil {
					return nil, fmt.Errorf("scan workspace: %w", err)
				}

				for _, file := range files {
					untrackedFiles.Add(file)
				}
			}

			continue
		}

		trackable, err := s.trackableFile(path, stat.FileInfo)
		if err != nil {
			return nil, fmt.Errorf("trackable file: %w", err)
		}

		if trackable {
			if stat.FileInfo.IsDir() {
				path += string(os.PathSeparator)
			} else {
				d := filepath.Dir(path) + string(os.PathSeparator)

				if untrackedFiles.Exists(d) {
					continue
				}
			}

			untrackedFiles.Add(path)
		}
	}

	return untrackedFiles.SortedValues(func(a, b string) bool {
		return a < b
	}), nil
}

func (s *StatusCommand) trackableFile(path string, stat os.FileInfo) (bool, error) {
	if stat == nil {
		return false, errors.New("stat is nil")
	}

	if stat.Mode().IsRegular() {
		return !s.repo.Index.Tracked(path), nil
	}

	if !stat.IsDir() {
		return false, nil
	}

	items, err := s.repo.Workspace.ListDir(path)
	if err != nil {
		return false, fmt.Errorf("list dir: %w", err)
	}

	for _, item := range items {
		trackable, err := s.trackableFile(filepath.Join(path, item.RelPath), item.FileInfo)
		if err != nil {
			return false, fmt.Errorf("trackable file: %w", err)
		}

		if trackable {
			return true, nil
		}
	}

	return false, nil
}

func (s *StatusCommand) Output(msg []byte, err error, stdout io.Writer) (int, error) {
	if err != nil {
		return 1, fmt.Errorf("output: %w", err)
	}

	fmt.Fprint(stdout, string(msg))

	return 0, nil
}
