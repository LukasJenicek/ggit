package repository

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/workspace"
)

// Repository
// Cwd = Is relative folder where you run ggit commands.
type Repository struct {
	FS        filesystem.Fs
	Workspace *workspace.Workspace
	Database  *database.Database

	Cwd         string
	RootDir     string
	GitPath     string
	Initialized bool
}

func New(fs filesystem.Fs, cwd string) (*Repository, error) {
	var initialized bool

	gitDir := filepath.Join(cwd, ".git")

	_, err := fs.Stat(gitDir)
	if err == nil {
		initialized = true
	}

	if err != nil && os.IsNotExist(err) {
		if os.IsNotExist(err) {
			initialized = false
		} else {
			return nil, fmt.Errorf("checking dir: %w", err)
		}
	}

	return &Repository{
		RootDir:     cwd,
		Cwd:         cwd,
		GitPath:     gitDir,
		Initialized: initialized,
		FS:          fs,
		Workspace:   workspace.New(cwd),
		Database:    database.New(gitDir),
	}, nil
}

func (r *Repository) Init() error {
	dirs := []string{".git", ".git/objects", ".git/refs"}

	for _, path := range dirs {
		err := r.FS.Mkdir(filepath.Join(r.RootDir, path), os.ModePerm)
		if err != nil {
			return fmt.Errorf("create %s directory: %w", path, err)
		}
	}

	return nil
}

func (r *Repository) Commit() error {
	files, err := r.Workspace.ListFiles()
	if err != nil {
		return fmt.Errorf("list files: %w", err)
	}

	blobs := make([]*database.Blob, 0, len(files))

	for _, file := range files {
		// TODO: Dirs not handled at the moment
		if file.Dir {
			continue
		}

		f, err := os.Open(file.Path)
		if err != nil {
			return fmt.Errorf("open file %s: %w", file.Path, err)
		}

		b, err := database.NewBlob(f, r.RootDir)
		if err != nil {
			return fmt.Errorf("create blob %s: %w", file.Path, err)
		}

		if err := r.Database.Store(b); err != nil {
			return fmt.Errorf("store blob %s: %w", file.Path, err)
		}

		blobs = append(blobs, b)
	}

	t := database.NewTree(blobs)
	if err := r.Database.Store(t); err != nil {
		return fmt.Errorf("store tree: %w", err)
	}

	now := time.Now()
	author := database.NewAuthor("lukas.jenicek5@gmail.com", "Lukas Jenicek", &now)

	c := database.NewCommit(hex.EncodeToString(t.ID()), author, "all")
	if err := r.Database.Store(c); err != nil {
		return fmt.Errorf("store commit: %w", err)
	}

	hFile, err := os.Create(r.GitPath + "/HEAD")
	if err != nil {
		return fmt.Errorf("create HEAD file: %w", err)
	}
	defer hFile.Close()

	cID := hex.EncodeToString(c.ID())
	if _, err := hFile.WriteString(cID); err != nil {
		return fmt.Errorf("write HEAD file: %w", err)
	}

	fmt.Println("commit successfully ", cID, " ", "all")

	return nil
}
