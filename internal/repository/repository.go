package repository

import (
	"fmt"
	"github.com/LukasJenicek/ggit/internal/clock"
	"github.com/LukasJenicek/ggit/internal/config"
	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/database/index"
	"github.com/LukasJenicek/ggit/internal/ds"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/workspace"
	"os"
	"path/filepath"
)

// Repository
// Cwd = Is relative folder where you run ggit commands.
type Repository struct {
	Config    *config.Config
	Clock     clock.Clock
	Database  *database.Database
	FS        filesystem.Fs
	Indexer   *index.Indexer
	Workspace *workspace.Workspace
	Refs      *database.Refs

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

	writer := filesystem.NewAtomicFileWriter(fs)

	refs, err := database.NewRefs(fs, gitPath, writer)
	if err != nil {
		return nil, fmt.Errorf("init refs: %w", err)
	}

	db := database.New(fs, gitPath)

	return &Repository{
		FS:          fs,
		Workspace:   workspace.New(cwd, fs),
		Database:    db,
		Indexer:     index.NewIndexer(fs, writer, filesystem.NewFileLocker(fs), db, gitPath, cwd),
		Clock:       clock,
		Refs:        refs,
		Config:      cfg,
		RootDir:     cwd,
		Cwd:         cwd,
		GitPath:     gitPath,
		Initialized: initialized,
	}, nil
}

func (r *Repository) Init() error {
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
		err := r.FS.Mkdir(filepath.Join(r.RootDir, path), os.ModePerm)
		if err != nil {
			return fmt.Errorf("create %s directory: %w", path, err)
		}
	}

	if err := r.Refs.UpdateHead("ref: refs/heads/master"); err != nil {
		return fmt.Errorf("update refs: %w", err)
	}

	return nil
}

func (r *Repository) Add(paths []string) error {
	files := []string{}

	for _, path := range paths {
		f, err := r.Workspace.ListFiles(path)
		if err != nil {
			return fmt.Errorf("list files: %w", err)
		}

		if path == "." {
			files = f

			break
		}

		files = append(files, f...)
	}

	if err := r.Indexer.Add(ds.NewSet(files)); err != nil {
		return fmt.Errorf("add files to index: %w", err)
	}

	return nil
}

func (r *Repository) Commit() (string, error) {
	files, err := r.Indexer.LoadIndex()
	if err != nil {
		return "", fmt.Errorf("load index: %w", err)
	}

	fmt.Println(files)

	return "", nil

	//files, err := r.Workspace.ListFiles(".")
	//if err != nil {
	//	return "", fmt.Errorf("list files: %w", err)
	//}
	//
	//entries, err := r.Database.SaveBlobs(ds.NewSet(files))
	//if err != nil {
	//	return "", fmt.Errorf("save blobs: %w", err)
	//}
	//
	//root, err := database.Build(database.NewTree(nil, ""), entries)
	//if err != nil {
	//	return "", fmt.Errorf("build tree: %w", err)
	//}
	//
	//rootID, err := r.Database.StoreTree(root)
	//if err != nil {
	//	return "", fmt.Errorf("store tree structure: %w", err)
	//}
	//
	//now := time.Now()
	//author := database.NewAuthor(r.Config.User.Email, r.Config.User.Name, &now)
	//
	//// TODO: read from file or cmd argument
	//commitMessage := "all"
	//
	//headCommit, err := r.Refs.Current()
	//if err != nil {
	//	return "", fmt.Errorf("get current commit: %w", err)
	//}
	//
	//c, err := database.NewCommit(hex.EncodeToString(rootID), author, commitMessage, headCommit)
	//if err != nil {
	//	return "", fmt.Errorf("create commit: %w", err)
	//}
	//
	//commitID, err := r.Database.Store(c)
	//if err != nil {
	//	return "", fmt.Errorf("store commit: %w", err)
	//}
	//
	//cID := hex.EncodeToString(commitID)
	//if err = r.Refs.UpdateHead(cID); err != nil {
	//	return "", fmt.Errorf("update head: %w", err)
	//}
	//
	//return cID, nil
}
