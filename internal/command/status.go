package command

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/LukasJenicek/ggit/internal/ds"
	"github.com/LukasJenicek/ggit/internal/index"
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
	index, err := s.repo.Index.Load()
	if err != nil {
		return nil, fmt.Errorf("load index entries: %w", err)
	}

	untrackedFiles, tracked, err := s.scanWorkspace(index, "")
	if err != nil {
		return nil, fmt.Errorf("scan workspace: %w", err)
	}

	spew.Dump(tracked)

	buf := bytes.NewBuffer(nil)
	for _, f := range untrackedFiles {
		buf.WriteString(fmt.Sprintf("?? %s\n", f))
	}

	return buf.Bytes(), nil
}

func (s *StatusCommand) scanWorkspace(index *index.Index, dirPrefix string) ([]string, map[string]fs.FileInfo, error) {
	stats, err := s.repo.Workspace.ListDir(dirPrefix)
	if err != nil {
		return nil, nil, fmt.Errorf("list dir: %w", err)
	}

	untrackedFiles := ds.Set[string]{}
	trackedFiles := make(map[string]fs.FileInfo)

	for _, stat := range stats {
		path := stat.RelPath

		if index.Tracked(path) {
			if !stat.FileInfo.IsDir() {
				trackedFiles[path] = stat.FileInfo
			}

			if stat.FileInfo.IsDir() {
				path = filepath.Join(dirPrefix, path)

				files, tracked, err := s.scanWorkspace(index, path)
				if err != nil {
					return nil, nil, fmt.Errorf("scan workspace: %w", err)
				}

				for _, file := range files {
					untrackedFiles.Add(file)
				}

				for p, trackedFile := range tracked {
					trackedFiles[p] = trackedFile
				}
			}

			continue
		}

		if dirPrefix != "" {
			path = filepath.Join(dirPrefix, path)
		}

		trackable, err := s.trackableFile(index, path, stat.FileInfo)
		if err != nil {
			return nil, nil, fmt.Errorf("trackable file: %w", err)
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
	}), trackedFiles, nil
}

func (s *StatusCommand) trackableFile(index *index.Index, path string, stat os.FileInfo) (bool, error) {
	if stat == nil {
		return false, errors.New("stat is nil")
	}

	if stat.Mode().IsRegular() {
		return !index.Tracked(path), nil
	}

	if !stat.IsDir() {
		return false, nil
	}

	items, err := s.repo.Workspace.ListDir(path)
	if err != nil {
		return false, fmt.Errorf("list dir: %w", err)
	}

	for _, item := range items {
		trackable, err := s.trackableFile(index, filepath.Join(path, item.RelPath), item.FileInfo)
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
