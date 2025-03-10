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

type Stat struct {
	RelPath  string
	FileInfo fs.FileInfo
}

// TODO: Should be relative to current working dir.
func (w Workspace) ListDir(dir string) ([]*Stat, error) {
	// TODO: load more ignored files from config
	ignore := []string{".", "..", ".git", ".idea"}

	path := w.rootDir
	if dir != "" {
		path = filepath.Join(w.rootDir, dir)
	}

	entries, err := w.fs.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("read dir %q: %w", dir, err)
	}

	stats := make([]*Stat, 0)

	for _, entry := range entries {
		if slices.Contains(ignore, entry.Name()) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("get entry info: %w", err)
		}

		stats = append(stats, &Stat{
			RelPath:  entry.Name(),
			FileInfo: info,
		})
	}

	return stats, nil
}

// ListFiles
// os.Walkdir: The files are walked in lexical order which makes the output deterministic.
func (w Workspace) ListFiles() ([]string, error) {
	// TODO: load more ignored files from config
	ignore := []string{".", "..", ".git", ".idea"}

	var files []string

	err := w.fs.WalkDir(w.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir: %w", err)
		}

		if slices.Contains(ignore, d.Name()) {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		cleanPath, err := filepath.Rel(w.rootDir, path)
		if err != nil {
			return fmt.Errorf("get relative path: %w", err)
		}

		files = append(files, cleanPath)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk dir: %w", err)
	}

	return files, nil
}

// MatchFiles
// Matches files against the provided pattern and returns their relative paths.
// It traverses the directory structure and applies pattern matching to find matching files.
func (w Workspace) MatchFiles(patternMatch string) ([]string, error) {
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
			if d.Name() == strings.TrimRight(patternMatch, "/") {
				matchDirectory = true
			}

			return nil
		}

		cleanPath, err := filepath.Rel(w.rootDir, path)
		if err != nil {
			return fmt.Errorf("get relative path: %w", err)
		}

		if patternMatch == "." {
			files = append(files, cleanPath)

			return nil
		}

		if matchDirectory && strings.Contains(cleanPath, patternMatch) {
			files = append(files, cleanPath)

			return nil
		}

		match, err := filepath.Match(patternMatch, cleanPath)
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
