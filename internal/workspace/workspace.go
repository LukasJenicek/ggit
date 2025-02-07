package workspace

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
)

type Workspace struct {
	Cwd string
}

func New(cwd string) *Workspace {
	return &Workspace{
		Cwd: cwd,
	}
}

type File struct {
	Path string
	Dir  bool
}

func (w Workspace) ListFiles() ([]*File, error) {
	// TODO: load more ignored files from config
	ignore := []string{".", "..", ".git"}

	files := []*File{}
	err := filepath.WalkDir(w.Cwd, func(path string, d fs.DirEntry, err error) error {
		if slices.Contains(ignore, path) {
			return nil
		}

		if slices.Contains(ignore, d.Name()) {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			files = append(files, &File{Path: path, Dir: false})
		} else {
			files = append(files, &File{Path: path, Dir: true})
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk recursively dir %q: %w", w.Cwd, err)
	}

	return files, nil
}
