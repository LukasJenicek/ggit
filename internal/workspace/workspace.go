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

type ErrPathNotMatched struct {
	Pattern string
}

func (e *ErrPathNotMatched) Error() string {
	return fmt.Sprintf("fatal: pathspec %q did not match any files", e.Pattern)
}

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

func (w Workspace) ListFiles(matchPath string) ([]string, error) {
	// TODO: load more ignored files from config
	ignore := []string{".", "..", ".git"}

	var files []string
	err := w.fs.WalkDir(w.rootDir, func(path string, d fs.DirEntry, err error) error {
		if slices.Contains(ignore, d.Name()) {
			return filepath.SkipDir
		}

		if err != nil {
			return err
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

	if len(files) == 0 {
		return nil, &ErrPathNotMatched{Pattern: matchPath}
	}

	return files, nil
}
