package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/LukasJenicek/ggit/internal/filesystem"
)

type Workspace struct {
	rootDir string
	fs      filesystem.Fs
}

func New(rootDir string, fs filesystem.Fs) *Workspace {
	return &Workspace{
		rootDir: rootDir,
		fs:      fs,
	}
}

type File struct {
	Path string
}

func (w Workspace) ListFiles(matchPath string) ([]string, error) {
	// TODO: load more ignored files from config
	ignore := []string{".", "..", ".git"}

	files := []string{}

	err := w.fs.WalkDir(w.rootDir, func(path string, d fs.DirEntry, err error) error {
		if slices.Contains(ignore, d.Name()) {
			return filepath.SkipDir
		}

		// do not include root folder name
		if filepath.Base(w.rootDir) == d.Name() {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		// prevent traversal directory attack
		cleanPath := filepath.Join(w.rootDir, filepath.Clean(strings.Replace(path, w.rootDir, "", 1)))
		if !strings.HasPrefix(cleanPath, w.rootDir) {
			return errors.New("invalid file path, outside of the root directory")
		}

		if matchPath == "." {
			files = append(files, cleanPath)
		}

		match, err := filepath.Match(matchPath, d.Name())
		if err != nil {
			return fmt.Errorf("matching path %q with pattern %q: %w", cleanPath, matchPath, err)
		}

		if match {
			files = append(files, cleanPath)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk recursively dir %q: %w", w.rootDir, err)
	}

	return files, nil
}
