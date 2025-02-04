package workspace

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
)

type Workspace struct{}

func New() *Workspace {
	return &Workspace{}
}

func (w Workspace) ListFiles(cwd string) ([]string, error) {
	// TODO: load more ignored files from config
	ignore := []string{".", "..", ".git"}

	files := []string{}

	err := filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
		if slices.Contains(ignore, path) {
			return nil
		}

		if slices.Contains(ignore, d.Name()) {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk recursively dir %q: %w", cwd, err)
	}

	return files, nil
}
