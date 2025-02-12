package repository_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/config"
	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/filesystem/memory"
	"github.com/LukasJenicek/ggit/internal/repository"
	"github.com/LukasJenicek/ggit/internal/workspace"
)

func TestNew(t *testing.T) {
	t.Parallel()

	cwd := "/home/test"
	gitPath := cwd + "/.git"

	t.Run("NotInitialized", func(t *testing.T) {
		t.Parallel()

		fs := memory.New()

		repo, err := repository.New(fs, cwd)
		require.NoError(t, err)

		refs, err := database.NewRefs(gitPath, &database.AtomicFileWriter{})
		require.NoError(t, err)

		require.NotNil(t, repo)
		require.EqualValues(t, &repository.Repository{
			FS:        fs,
			Workspace: workspace.New(cwd),
			Database:  database.New(gitPath),
			Refs:      refs,
			Config: &config.Config{
				User: &config.User{
					Name:  "Lukas Jenicek",
					Email: "lukas.jenicek5@gmail.com",
				},
			},
			Cwd:         cwd,
			RootDir:     cwd,
			GitPath:     gitPath,
			Initialized: false,
		}, repo)
	})

	t.Run("AlreadyInitialized", func(t *testing.T) {
		t.Parallel()

		fs := memory.New()
		err := fs.Mkdir("/home/test/.git", os.ModePerm)
		require.NoError(t, err)

		repo, err := repository.New(fs, cwd)
		require.NoError(t, err)
		require.NotNil(t, repo)

		refs, err := database.NewRefs(gitPath, &database.AtomicFileWriter{})
		require.NoError(t, err)

		require.EqualValues(t, &repository.Repository{
			FS:        fs,
			Workspace: workspace.New(cwd),
			Database:  database.New(gitPath),
			Refs:      refs,
			Config: &config.Config{
				User: &config.User{
					Name:  "Lukas Jenicek",
					Email: "lukas.jenicek5@gmail.com",
				},
			},
			Cwd:         cwd,
			RootDir:     cwd,
			GitPath:     gitPath,
			Initialized: true,
		}, repo)
	})
}

func TestInit(t *testing.T) {
	t.Parallel()

	t.Run("Init", func(t *testing.T) {
		t.Parallel()

		repo, err := repository.New(memory.New(), "/home/test/")
		require.NoError(t, err)
		require.NotNil(t, repo)

		err = repo.Init()
		require.NoError(t, err)

		g, err := repo.FS.Stat("/home/test/.git")
		require.NoError(t, err)
		require.True(t, g.IsDir())

		g, err = repo.FS.Stat("/home/test/.git/refs")
		require.NoError(t, err)
		require.True(t, g.IsDir())

		g, err = repo.FS.Stat("/home/test/.git/objects")
		require.NoError(t, err)
		require.True(t, g.IsDir())
	})
}
