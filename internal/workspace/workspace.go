package workspace

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/LukasJenicek/ggit/internal/filesystem"
)

type Workspace struct {
	Cwd string
	fs  filesystem.Fs
}

func New(cwd string, fs filesystem.Fs) *Workspace {
	return &Workspace{
		Cwd: cwd,
		fs:  fs,
	}
}

type File struct {
	Path string
	Dir  bool
}

func (w Workspace) ListFiles(matchPath string) ([]*File, error) {
	// TODO: load more ignored files from config
	ignore := []string{".", "..", ".git"}

	files := []*File{}

	err := w.fs.WalkDir(w.Cwd, func(path string, d fs.DirEntry, err error) error {
		if slices.Contains(ignore, d.Name()) {
			return filepath.SkipDir
		}

		// do not include root folder name
		if filepath.Base(w.Cwd) == d.Name() {
			return nil
		}

		if !d.IsDir() {
			if matchPath == "." {
				files = append(files, &File{Path: path, Dir: true})
			}

			match, err := filepath.Match(matchPath, d.Name())
			if err != nil {
				return fmt.Errorf("matching path %q: %w", matchPath, err)
			}

			if match {
				files = append(files, &File{Path: path, Dir: false})
			}
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
