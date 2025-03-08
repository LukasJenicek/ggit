package repository

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/LukasJenicek/ggit/internal/clock"
	"github.com/LukasJenicek/ggit/internal/config"
	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/index"
	"github.com/LukasJenicek/ggit/internal/workspace"
)

var ErrNoFilesToCommit = errors.New("nothing added to commit (use 'ggit add' to track)")

// Repository
// Cwd = Is relative folder where you run ggit commands.
type Repository struct {
	GitConfig *config.Config
	Database  *database.Database
	Index     *index.Indexer
	Refs      *database.Refs
	Workspace *workspace.Workspace

	// Important for unit testing
	Clock clock.Clock
	FS    filesystem.Fs

	Cwd         string
	RootDir     string
	GitPath     string
	Initialized bool
}

func New(fs filesystem.Fs, clock clock.Clock, cwd string) (*Repository, error) {
	var initialized bool

	gitPath := filepath.Join(cwd, ".git")

	_, err := fs.Stat(gitPath)
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

	cfg, err := config.LoadGitConfig()
	if err != nil {
		return nil, fmt.Errorf("load git config: %w", err)
	}

	locker := filesystem.NewFileLocker(fs)

	writer, err := filesystem.NewAtomicFileWriter(fs, locker)
	if err != nil {
		return nil, fmt.Errorf("create file writer: %w", err)
	}

	refs, err := database.NewRefs(fs, gitPath, writer)
	if err != nil {
		return nil, fmt.Errorf("init refs: %w", err)
	}

	db, err := database.New(fs, cwd)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	indexer, err := index.NewIndexer(fs, writer, locker, db, cwd)
	if err != nil {
		return nil, fmt.Errorf("init indexer: %w", err)
	}

	w, err := workspace.New(cwd, fs)
	if err != nil {
		return nil, fmt.Errorf("init workspace: %w", err)
	}

	return &Repository{
		FS:          fs,
		Workspace:   w,
		Database:    db,
		Index:       indexer,
		Clock:       clock,
		Refs:        refs,
		GitConfig:   cfg,
		RootDir:     cwd,
		Cwd:         cwd,
		GitPath:     gitPath,
		Initialized: initialized,
	}, nil
}

func (repo *Repository) Init() error {
	dirs := []string{
		".git",
		".git/branches",
		".git/hooks",
		".git/info",
		".git/objects",
		".git/objects/info",
		".git/objects/pack",
		".git/refs",
		".git/refs/heads",
		".git/refs/tags",
	}

	for _, path := range dirs {
		absPath := filepath.Join(repo.RootDir, path)

		err := repo.FS.Mkdir(absPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("create %s directory: %w", absPath, err)
		}
	}

	// TODO: Parse default branch from config
	if err := repo.Refs.InitRef("ref: refs/heads/master"); err != nil {
		return fmt.Errorf("init ref: %w", err)
	}

	return nil
}

func (repo *Repository) Add(paths []string) error {
	var files []string

	for _, path := range paths {
		f, err := repo.Workspace.MatchFiles(path)
		if err != nil {
			return fmt.Errorf("match files: %w", err)
		}

		if path == "." {
			files = f

			break
		}

		files = append(files, f...)
	}

	if err := repo.Index.Add(files); err != nil {
		return fmt.Errorf("add files to index: %w", err)
	}

	return nil
}

func (repo *Repository) Commit() (*database.Commit, error) {
	idx, err := repo.Index.Load()
	if err != nil {
		return nil, fmt.Errorf("load index: %w", err)
	}

	entriesLen := idx.Entries.Len()
	if entriesLen == 0 {
		return nil, ErrNoFilesToCommit
	}

	entries := make([]*database.Entry, 0, entriesLen)

	for _, iEntry := range idx.Entries.SortedValues() {
		filePath := string(iEntry.Path)

		entry, err := database.NewEntry(filepath.Base(filePath), filePath, iEntry.OID, false)
		if err != nil {
			return nil, fmt.Errorf("create entry: %w", err)
		}

		entries = append(entries, entry)
	}

	root, err := database.Build(database.NewRootTree(), entries)
	if err != nil {
		return nil, fmt.Errorf("build tree: %w", err)
	}

	rootID, err := repo.Database.StoreTree(root)
	if err != nil {
		return nil, fmt.Errorf("store tree structure: %w", err)
	}

	now := time.Now()
	author := database.NewAuthor(repo.GitConfig.User.Email, repo.GitConfig.User.Name, &now)

	// TODO: read from file or cmd argument
	commitMessage := "all"

	parent, err := repo.Refs.ReadHead()
	if err != nil {
		return nil, fmt.Errorf("read head: %w", err)
	}

	c, err := database.NewCommit(parent, hex.EncodeToString(rootID), author, commitMessage)
	if err != nil {
		return nil, fmt.Errorf("create commit: %w", err)
	}

	commitID, err := repo.Database.Store(c)
	if err != nil {
		return nil, fmt.Errorf("store commit: %w", err)
	}

	cID := hex.EncodeToString(commitID)
	if err = repo.Refs.UpdateHead(cID); err != nil {
		return nil, fmt.Errorf("update head: %w", err)
	}

	if err = c.SetOID(cID); err != nil {
		return nil, fmt.Errorf("set commit oid: %w", err)
	}

	return c, nil
}
