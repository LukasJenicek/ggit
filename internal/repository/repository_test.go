package repository_test

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"syscall"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/clock"
	"github.com/LukasJenicek/ggit/internal/config"
	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/filesystem/memory"
	"github.com/LukasJenicek/ggit/internal/index"
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

		locker := filesystem.NewFileLocker(fs)

		writer, err := filesystem.NewAtomicFileWriter(fs, locker)
		require.NoError(t, err)

		refs, err := database.NewRefs(fs, gitPath, writer)
		require.NoError(t, err)

		require.NotNil(t, repo)

		d, err := database.New(fs, cwd)
		require.NoError(t, err)

		w, err := workspace.New(cwd, fs)
		require.NoError(t, err)

		indexer, err := index.NewIndexer(fs, writer, locker, d, cwd)
		require.NoError(t, err)

		require.EqualValues(t, &repository.Repository{
			FS:        fs,
			Workspace: w,
			Database:  d,
			Index:     indexer,
			Clock:     fakeClock,
			Refs:      refs,
			GitConfig: &config.Config{
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

		locker := filesystem.NewFileLocker(fs)

		writer, err := filesystem.NewAtomicFileWriter(fs, locker)
		require.NoError(t, err)

		refs, err := database.NewRefs(fs, gitPath, writer)
		require.NoError(t, err)

		require.NotNil(t, repo)

		db, err := database.New(fs, cwd)
		require.NoError(t, err)

		w, err := workspace.New(cwd, fs)
		require.NoError(t, err)

		indexer, err := index.NewIndexer(fs, writer, locker, db, cwd)
		require.NoError(t, err)

		require.EqualValues(t, &repository.Repository{
			FS:        fs,
			Workspace: w,
			Database:  db,
			Index:     indexer,
			Clock:     fakeClock,
			Refs:      refs,
			GitConfig: &config.Config{
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

func TestRepository_Commit(t *testing.T) {
	t.Parallel()

	cwd := "tmp/test"

	fs := memory.New(fstest.MapFS{
		"tmp/test/hello.txt": &fstest.MapFile{
			Data: []byte("hello"),
			Sys: &syscall.Stat_t{
				Dev:     66306,
				Ino:     26874043,
				Nlink:   1,
				Mode:    33204,
				Uid:     1000,
				Gid:     1000,
				X__pad0: 0,
				Rdev:    0,
				Size:    7,
				Blksize: 4096,
				Blocks:  8,
				Atim: syscall.Timespec{
					Sec:  1740497047,
					Nsec: 592315384,
				},
				Mtim: syscall.Timespec{
					Sec:  1739287401,
					Nsec: 888108884,
				},
				Ctim: syscall.Timespec{
					Sec:  1739287401,
					Nsec: 888108884,
				},
			},
		},
		"tmp/test/world.txt": &fstest.MapFile{
			Data: []byte("world"),
			Sys: &syscall.Stat_t{
				Dev:     66306,
				Ino:     26874043,
				Nlink:   1,
				Mode:    33204,
				Uid:     1000,
				Gid:     1000,
				Size:    6,
				Blksize: 4096,
				Blocks:  8,
				Atim: syscall.Timespec{
					Sec:  1740497047,
					Nsec: 592315384,
				},
				Mtim: syscall.Timespec{
					Sec:  1739287401,
					Nsec: 888108884,
				},
				Ctim: syscall.Timespec{
					Sec:  1739287401,
					Nsec: 888108884,
				},
			},
		},
	})

	now := time.Date(2024, 12, 15, 17, 8, 0o0, 0, time.UTC)

	repo, err := repository.New(
		fs,
		clock.NewFakeClock(now),
		cwd,
	)
	require.NoError(t, err)

	err = repo.Init()
	require.NoError(t, err)

	err = repo.Add([]string{"hello.txt", "world.txt"})
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
