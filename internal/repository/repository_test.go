package repository_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/libs/filesystem/memory"
	"github.com/LukasJenicek/ggit/libs/repository"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("NotInitialized", func(t *testing.T) {
		fs := memory.New()

		repo, err := repository.New(fs, "/home/test/")
		require.NoError(t, err)

		require.NotNil(t, repo)
		require.EqualValues(t, &repository.Repository{Cwd: "/home/test/", RootDir: "/home/test/", Initialized: false, FS: fs}, repo)
	})

	t.Run("AlreadyInitialized", func(t *testing.T) {
		fs := memory.New()
		err := fs.Mkdir("/home/test/.git", os.ModePerm)
		require.NoError(t, err)

		repo, err := repository.New(fs, "/home/test/")
		require.NoError(t, err)

		require.NotNil(t, repo)
		require.EqualValues(t, &repository.Repository{Cwd: "/home/test/", RootDir: "/home/test/", Initialized: true, FS: fs}, repo)
	})
}

func TestInit(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
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
