package index_test

import (
	"os"
	"syscall"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/filesystem/memory"
	"github.com/LukasJenicek/ggit/internal/index"
)

func TestAddEntryToIndex(t *testing.T) {
	t.Parallel()

	rootDir := "tmp/test"

	fs := memory.New(fstest.MapFS{
		"tmp/test": &fstest.MapFile{
			Mode: os.ModeDir,
			Sys:  defaultStat(uint32(os.ModeDir), 0),
		},
		"tmp/test/hello.txt": &fstest.MapFile{
			Data: []byte("hello"),
			Mode: 0o644,
			Sys:  defaultStat(uint32(0o644), 6),
		},
		"tmp/test/world.txt": &fstest.MapFile{
			Data: []byte("world"),
			Sys:  defaultStat(uint32(0o644), 6),
		},
	})
	locker := filesystem.NewFileLocker(fs)
	fileWriter, err := filesystem.NewAtomicFileWriter(fs, locker)
	require.NoError(t, err)

	db, err := database.New(fs, rootDir)
	require.NoError(t, err)

	indexer, err := index.NewIndexer(fs, fileWriter, locker, db, rootDir)
	require.NoError(t, err)

	err = indexer.Add([]string{"hello.txt", "world.txt"})
	require.NoError(t, err)

	idx, err := indexer.Load()
	require.NoError(t, err)

	entriesNames := make([]string, idx.Entries.Len())
	for i, entry := range idx.Entries.SortedValues() {
		entriesNames[i] = string(entry.Path)
	}

	require.EqualValues(t, []string{"hello.txt", "world.txt"}, entriesNames)
}

func TestReplacesFileWithDirectory(t *testing.T) {
	t.Parallel()

	rootDir := "tmp/test"
	mapFS := fstest.MapFS{
		"tmp/test": &fstest.MapFile{
			Mode: os.ModeDir,
			Sys:  defaultStat(uint32(os.ModeDir), 0),
		},
		"tmp/test/hello.txt": &fstest.MapFile{
			Data: []byte("hello"),
			Mode: 0o644,
			Sys:  defaultStat(uint32(0o644), 6),
		},
		"tmp/test/world.txt": &fstest.MapFile{
			Data: []byte("world"),
			Sys:  defaultStat(uint32(0o644), 6),
		},
	}

	fs := memory.New(mapFS)
	locker := filesystem.NewFileLocker(fs)
	fileWriter, err := filesystem.NewAtomicFileWriter(fs, locker)
	require.NoError(t, err)

	db, err := database.New(fs, rootDir)
	require.NoError(t, err)

	idx, err := index.NewIndexer(fs, fileWriter, locker, db, rootDir)
	require.NoError(t, err)

	err = idx.Add([]string{"hello.txt", "world.txt"})
	require.NoError(t, err)

	mapFS["tmp/test/world.txt"] = &fstest.MapFile{
		Mode: os.ModeDir,
		Sys:  defaultStat(uint32(os.ModeDir), 0),
	}

	mapFS["tmp/test/world.txt/world.txt"] = &fstest.MapFile{
		Data: []byte("world"),
		Sys:  defaultStat(uint32(0o644), 6),
	}

	err = idx.Add([]string{"hello.txt", "world.txt/world.txt"})
	require.NoError(t, err)

	index, err := idx.Load()
	require.NoError(t, err)

	entriesNames := make([]string, index.Entries.Len())
	for i, entry := range index.Entries.SortedValues() {
		entriesNames[i] = string(entry.Path)
	}

	require.EqualValues(t, []string{"hello.txt", "world.txt/world.txt"}, entriesNames)
}

func TestReplacesDirectoryWithFile(t *testing.T) {
	t.Parallel()

	rootDir := "tmp/test"
	mapFS := fstest.MapFS{
		"tmp/test": &fstest.MapFile{
			Mode: os.ModeDir,
			Sys:  defaultStat(uint32(os.ModeDir), 0),
		},
		"tmp/test/hello.txt": &fstest.MapFile{
			Mode: os.ModeDir,
			Sys:  defaultStat(uint32(os.ModeDir), 0),
		},
		"tmp/test/hello.txt/hello.txt": &fstest.MapFile{
			Data: []byte("hello"),
			Mode: 0o644,
			Sys:  defaultStat(uint32(0o644), 6),
		},
		"tmp/test/world.txt": &fstest.MapFile{
			Data: []byte("world"),
			Sys:  defaultStat(uint32(0o644), 6),
		},
	}

	fs := memory.New(mapFS)
	locker := filesystem.NewFileLocker(fs)
	fileWriter, err := filesystem.NewAtomicFileWriter(fs, locker)
	require.NoError(t, err)

	db, err := database.New(fs, rootDir)
	require.NoError(t, err)

	idx, err := index.NewIndexer(fs, fileWriter, locker, db, rootDir)
	require.NoError(t, err)

	err = idx.Add([]string{"hello.txt/hello.txt", "world.txt"})
	require.NoError(t, err)

	mapFS["tmp/test/hello.txt"] = &fstest.MapFile{
		Data: []byte("hello"),
		Sys:  defaultStat(uint32(0o644), 6),
	}

	err = idx.Add([]string{"hello.txt"})
	require.NoError(t, err)

	index, err := idx.Load()
	require.NoError(t, err)

	entriesNames := make([]string, index.Entries.Len())
	for i, entry := range index.Entries.SortedValues() {
		entriesNames[i] = string(entry.Path)
	}

	require.EqualValues(t, []string{"hello.txt", "world.txt"}, entriesNames)
}

func defaultStat(mode uint32, size int64) *syscall.Stat_t {
	return &syscall.Stat_t{
		Dev:     66306,
		Ino:     26874043,
		Nlink:   1,
		Mode:    mode,
		Uid:     1000,
		Gid:     1000,
		Size:    size,
		Blksize: 4096,
		Blocks:  (size + 511) / 512, // Approximate block count
		Atim:    syscall.Timespec{Sec: 1740497047, Nsec: 592315384},
		Mtim:    syscall.Timespec{Sec: 1739287401, Nsec: 888108884},
		Ctim:    syscall.Timespec{Sec: 1739287401, Nsec: 888108884},
	}
}
