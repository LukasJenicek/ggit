package repository_test

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/clock"
	"github.com/LukasJenicek/ggit/internal/config"
	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/database/index"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/filesystem/memory"
	"github.com/LukasJenicek/ggit/internal/repository"
	"github.com/LukasJenicek/ggit/internal/workspace"
)

func TestNew(t *testing.T) {
	t.Parallel()

	cwd := "home/test"
	gitPath := cwd + "/.git"

	fakeClock := clock.NewFakeClock(time.Date(2000, 12, 15, 17, 8, 0o0, 0, time.UTC))

	t.Run("NotInitialized", func(t *testing.T) {
		t.Parallel()

		fs := memory.New(fstest.MapFS{})

		repo, err := repository.New(
			fs,
			fakeClock,
			cwd,
		)
		require.NoError(t, err)

		refs, err := database.NewRefs(fs, gitPath, filesystem.NewAtomicFileWriter(fs))
		require.NoError(t, err)

		require.NotNil(t, repo)

		d := database.New(fs, gitPath)

		require.EqualValues(t, &repository.Repository{
			FS:        fs,
			Workspace: workspace.New(cwd, fs),
			Database:  d,
			Indexer:   index.NewIndexer(fs, filesystem.NewAtomicFileWriter(fs), filesystem.NewFileLocker(fs), d, gitPath, cwd),
			Clock:     fakeClock,
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

		fs := memory.New(fstest.MapFS{})
		err := fs.Mkdir("home/test/.git", os.ModePerm)
		require.NoError(t, err)

		repo, err := repository.New(fs, fakeClock, cwd)
		require.NoError(t, err)
		require.NotNil(t, repo)

		writer := filesystem.NewAtomicFileWriter(fs)

		refs, err := database.NewRefs(fs, gitPath, writer)
		require.NoError(t, err)

		db := database.New(fs, gitPath)

		require.EqualValues(t, &repository.Repository{
			FS:        fs,
			Workspace: workspace.New(cwd, fs),
			Database:  db,
			Indexer:   index.NewIndexer(fs, writer, filesystem.NewFileLocker(fs), db, gitPath, cwd),
			Clock:     fakeClock,
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

		repo, err := repository.New(
			memory.New(fstest.MapFS{}),
			clock.NewFakeClock(time.Date(2000, 12, 15, 17, 8, 0o0, 0, time.UTC)),
			"tmp/test/",
		)
		require.NoError(t, err)
		require.NotNil(t, repo)

		err = repo.Init()
		require.NoError(t, err)

		g, err := repo.FS.Stat("tmp/test/.git")
		require.NoError(t, err)
		require.True(t, g.IsDir())

		g, err = repo.FS.Stat("tmp/test/.git/refs")
		require.NoError(t, err)
		require.True(t, g.IsDir())

		g, err = repo.FS.Stat("tmp/test/.git/objects")
		require.NoError(t, err)
		require.True(t, g.IsDir())
	})
}

func TestRepository_Commit(t *testing.T) {
	t.Skip()
	t.Parallel()

	cwd := "tmp/test"

	fs := memory.New(fstest.MapFS{
		"tmp/test/hello.txt": &fstest.MapFile{
			Data: []byte("hello"),
		},
		"tmp/test/world.txt": &fstest.MapFile{
			Data: []byte("world"),
		},
	})

	now := time.Date(2024, 12, 15, 17, 8, 0o0, 0, time.UTC)

	repo, err := repository.New(
		fs,
		clock.NewFakeClock(now),
		cwd,
	)
	require.NoError(t, err)

	_, err = repo.Commit()
	require.NoError(t, err)

	helloBlob := hash(t, database.NewBlob([]byte("hello")))
	worldBlob := hash(t, database.NewBlob([]byte("world")))

	root := database.NewTree(nil, "")

	entry, err := database.NewEntry("hello.txt", "tmp/test/hello.txt", helloBlob, false)
	require.NoError(t, err)
	root.AddEntry(entry)

	entry, err = database.NewEntry("world.txt", "tmp/test/world.txt", worldBlob, false)
	require.NoError(t, err)
	root.AddEntry(entry)

	gitFileObjects := []string{
		hex.EncodeToString(helloBlob),
		hex.EncodeToString(worldBlob),
		hex.EncodeToString(hash(t, root)),
	}

	for _, file := range gitFileObjects {
		_, err = fs.Open("tmp/test/.git/objects/" + file[0:2] + "/" + file[2:])
		require.NoError(t, err)
	}
}

func hash(t *testing.T, object database.Object) []byte {
	t.Helper()

	content, err := object.Content()
	require.NoError(t, err)

	hasher := sha1.New()
	hasher.Write(content)

	return hasher.Sum(nil)
}
