package command_test

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/clock"
	"github.com/LukasJenicek/ggit/internal/command"
	"github.com/LukasJenicek/ggit/internal/filesystem/memory"
	"github.com/LukasJenicek/ggit/internal/hasher"
	"github.com/LukasJenicek/ggit/internal/index"
	"github.com/LukasJenicek/ggit/internal/repository"
)

func TestAddFilesNameConflict(t *testing.T) {
	t.Parallel()

	fs := fstest.MapFS{
		"tmp/test/": &fstest.MapFile{
			Mode: os.ModeDir,
		},
		"tmp/test/world.txt": &fstest.MapFile{
			Data: []byte("world"),
			Mode: 33204,
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
	}

	repo, err := repository.New(
		memory.New(fs),
		clock.NewFakeClock(time.Date(2000, 12, 15, 17, 8, 0o0, 0, time.UTC)),
		"tmp/test",
	)
	require.NoError(t, err)

	initCmd, err := command.NewInitCommand(repo)
	require.NoError(t, err)

	_, err = initCmd.Run()
	require.NoError(t, err)

	cmd, err := command.NewAddCommand([]string{"."}, repo)
	require.NoError(t, err)

	_, err = cmd.Run()
	require.NoError(t, err)

	fs["tmp/test/hello.txt"] = &fstest.MapFile{
		Mode: os.ModeDir,
	}

	fs["tmp/test/hello.txt/world.txt"] = &fstest.MapFile{
		Data: []byte("world"),
		Mode: 33204,
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
	}

	cmd, err = command.NewAddCommand([]string{"."}, repo)
	require.NoError(t, err)

	_, err = cmd.Run()
	require.NoError(t, err)

	content, err := fs.ReadFile("tmp/test/.git/index")
	require.NoError(t, err)

	var expectedContent []byte
	expectedContent = append(expectedContent, []byte("DIRC")...)
	// version
	expectedContent = append(expectedContent, []byte{0, 0, 0, 2}...)
	// number of files
	expectedContent = append(expectedContent, []byte{0, 0, 0, 2}...)

	bytes, err := fileContent(t, fs, "tmp/test/hello.txt/world.txt", "hello.txt/world.txt")
	require.NoError(t, err)

	expectedContent = append(expectedContent, bytes...)

	bytes, err = fileContent(t, fs, "tmp/test/world.txt", "world.txt")
	require.NoError(t, err)

	expectedContent = append(expectedContent, bytes...)

	// sha1 checksum of content
	oid, err := hasher.SHA1HashContent(expectedContent)
	require.NoError(t, err)

	expectedContent = append(expectedContent, oid...)

	require.EqualValues(t, expectedContent, content)
}

func TestAddFiles(t *testing.T) {
	t.Parallel()

	fs := fstest.MapFS{
		"tmp/test/": &fstest.MapFile{
			Mode: os.ModeDir,
		},
		"tmp/test/hello.txt": &fstest.MapFile{
			Data: []byte("hello"),
			Mode: 0o644,
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
			Mode: 33204,
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
	}

	repo, err := repository.New(
		memory.New(fs),
		clock.NewFakeClock(time.Date(2000, 12, 15, 17, 8, 0o0, 0, time.UTC)),
		"tmp/test",
	)
	require.NoError(t, err)

	initCmd, err := command.NewInitCommand(repo)
	require.NoError(t, err)

	_, err = initCmd.Run()
	require.NoError(t, err)

	cmd, err := command.NewAddCommand([]string{"."}, repo)
	require.NoError(t, err)

	_, err = cmd.Run()
	require.NoError(t, err)

	content, err := fs.ReadFile("tmp/test/.git/index")
	require.NoError(t, err)

	var expectedContent []byte
	expectedContent = append(expectedContent, []byte("DIRC")...)
	// version
	expectedContent = append(expectedContent, []byte{0, 0, 0, 2}...)
	// number of files
	expectedContent = append(expectedContent, []byte{0, 0, 0, 2}...)

	bytes, err := fileContent(t, fs, "tmp/test/hello.txt", "hello.txt")
	require.NoError(t, err)

	expectedContent = append(expectedContent, bytes...)

	bytes, err = fileContent(t, fs, "tmp/test/world.txt", "world.txt")
	require.NoError(t, err)

	expectedContent = append(expectedContent, bytes...)

	// sha1 checksum of content
	oid, err := hasher.SHA1HashContent(expectedContent)
	require.NoError(t, err)

	expectedContent = append(expectedContent, oid...)

	require.EqualValues(t, expectedContent, content)
}

func fileContent(t *testing.T, fs fstest.MapFS, filepath, filename string) ([]byte, error) {
	t.Helper()

	stat, err := fs.Stat(filepath)
	require.NoError(t, err)

	fContent, err := fs.ReadFile(filepath)
	require.NoError(t, err)

	oid, err := hasher.SHA1HashContent([]byte(fmt.Sprintf("%s %d\x00%s", "blob", len(fContent), fContent)))
	require.NoError(t, err)

	entry, err := index.NewEntry(filename, stat, oid)
	require.NoError(t, err)

	content, err := entry.Content()
	require.NoError(t, err)

	return content, nil
}
