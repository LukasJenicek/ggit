package repository

import (
	"fmt"
	"github.com/LukasJenicek/ggit/internal/config"
	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/workspace"
	"github.com/davecgh/go-spew/spew"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Repository
// Cwd = Is relative folder where you run ggit commands.
type Repository struct {
	FS        filesystem.Fs
	Workspace *workspace.Workspace
	Database  *database.Database
	Refs      *database.Refs
	Config    *config.Config

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

	cfg, err := config.LoadGitConfig()
	if err != nil {
		return nil, fmt.Errorf("load git config: %w", err)
	}

	refs, err := database.NewRefs(gitDir, &database.AtomicFileWriter{})
	if err != nil {
		return nil, fmt.Errorf("init refs: %w", err)
	}

	return &Repository{
		FS:          fs,
		Workspace:   workspace.New(cwd),
		Database:    database.New(gitDir),
		Refs:        refs,
		Config:      cfg,
		RootDir:     cwd,
		Cwd:         cwd,
		GitPath:     gitDir,
		Initialized: initialized,
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

	entries := make([]*database.Entry, 0, len(files))
	for _, file := range files {
		if file.Dir {
			continue
		}

		f, err := os.Open(file.Path)
		if err != nil {
			return fmt.Errorf("open file %s: %w", file.Path, err)
		}

		content, err := os.ReadFile(f.Name())
		if err != nil {
			return fmt.Errorf("read file content: %w", err)
		}

		relativeFilePath := strings.Replace(f.Name(), r.RootDir+"/", "", 1)

		i, err := os.Stat(f.Name())
		if err != nil {
			return fmt.Errorf("stat file %s: %w", f.Name(), err)
		}

		executable := false
		if i.Mode().Perm()&0o100 != 0 {
			executable = true
		}

		b, err := database.NewBlob(content)
		if err != nil {
			return fmt.Errorf("create blob %s: %w", file.Path, err)
		}

		oid, err := r.Database.Store(b)
		if err != nil {
			return fmt.Errorf("store blob %s: %w", file.Path, err)
		}

		entries = append(entries, database.NewEntry(relativeFilePath, oid, executable))
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	spew.Dump(entries)

	entriesMap := make(map[string]database.Object)
	for _, entry := range entries {
		d := filepath.Dir(entry.Name)

		// root folder
		if d == "." {
			entriesMap[entry.Name] = entry
			continue
		}

		t, ok := entriesMap[d]
		if !ok {
			t = database.NewTree()
			entriesMap[d] = t
		}

		tree, ok := t.(*database.Tree)
		if !ok {
			return fmt.Errorf("unexpected type %T", t)
		}

		tree.AddEntry(entry)
	}

	return nil
	//

	//root := database.NewTree(entries)
	//if err := r.Database.Store(root); err != nil {
	//	return fmt.Errorf("store tree: %w", err)
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
	//	return fmt.Errorf("get current commit: %w", err)
	//}
	//
	//c, err := database.NewCommit(hex.EncodeToString(root.ID()), author, commitMessage, headCommit)
	//if err != nil {
	//	return fmt.Errorf("create commit: %w", err)
	//}
	//
	//if err := r.Database.Store(c); err != nil {
	//	return fmt.Errorf("store commit: %w", err)
	//}
	//
	//cID := hex.EncodeToString(c.ID())
	//if err = r.Refs.UpdateHead(cID); err != nil {
	//	return fmt.Errorf("update head: %w", err)
	//}
	//
	//fmt.Println("commit successfully ", cID, " ", commitMessage)
	//
	//return nil
}
