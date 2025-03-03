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

func New(rootDir string, fs filesystem.Fs) (*Workspace, error) {
	if fs == nil {
		return nil, errors.New("filesystem must not be nil")
	}

	if rootDir == "" {
		return nil, errors.New("root dir must not be empty")
	}

	return &Workspace{
		rootDir: rootDir,
		fs:      fs,
	}, nil
}

func (w Workspace) ListFiles(patternMatch string) ([]string, error) {
	// TODO: load more ignored files from config
	ignore := []string{".", "..", ".git"}

	var files []string

	matchDirectory := false

	err := w.fs.WalkDir(w.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir: %w", err)
		}

		if slices.Contains(ignore, d.Name()) {
			return filepath.SkipDir
		}

		// do not include root folder name
		if filepath.Base(w.rootDir) == d.Name() {
			return nil
		}

		if d.IsDir() {
			if d.IsDir() && d.Name() == strings.TrimRight(patternMatch, "/") {
				matchDirectory = true
			}

			return nil
		}

		cleanPath, err := filepath.Rel(w.rootDir, filepath.Join(w.rootDir, path))
		if err != nil {
			return fmt.Errorf("filepath.Rel(%q, %q): %w", w.rootDir, path, err)
		}

		if patternMatch == "." {
			files = append(files, cleanPath)

			return nil
		}

		if matchDirectory && strings.Contains(cleanPath, patternMatch) {
			files = append(files, cleanPath)

			return nil
		}

		match, err := filepath.Match(filepath.Join(w.rootDir, patternMatch), cleanPath)
		if err != nil {
			return fmt.Errorf("matching path %q with pattern %q: %w", cleanPath, patternMatch, err)
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
		return nil, &ErrPathNotMatched{Pattern: patternMatch}
	}

	return files, nil
}
