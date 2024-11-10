package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LukasJenicek/ggit/libs/filesystem"
)

// Repository
// RootDir = Where .git folder exist
// Cwd = Is relative folder where you run ggit commands
type Repository struct {
	Cwd         string
	RootDir     string
	Initialized bool
	FS          filesystem.Fs
}

func New(fs filesystem.Fs, cwd string) (*Repository, error) {
	var initialized bool

	gitDir := filepath.Join(cwd, ".git")
	_, err := fs.Stat(gitDir)
	if err == nil {
		initialized = true
	}

	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("checking dir: %w", err)
	}

	return &Repository{
		RootDir:     cwd,
		Cwd:         cwd,
		Initialized: initialized,
		FS:          fs,
	}, nil
}

func (r *Repository) Init() error {
	dirs := []string{".git", ".git/objects", ".git/refs"}

	for _, path := range dirs {
		err := os.Mkdir(filepath.Join(r.RootDir, path), os.ModePerm)

		if err != nil {
			return fmt.Errorf("create %s directory: %r", path, err)
		}
	}

	return nil
}
