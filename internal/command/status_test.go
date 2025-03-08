package command_test

import (
	"os"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/clock"
	"github.com/LukasJenicek/ggit/internal/command"
	"github.com/LukasJenicek/ggit/internal/filesystem/memory"
	"github.com/LukasJenicek/ggit/internal/repository"
)

func TestListingUntrackedFiles(t *testing.T) {
	t.Parallel()

	repo, err := repository.New(
		memory.New(fstest.MapFS{
			"tmp/test/": &fstest.MapFile{
				Mode: os.ModeDir,
			},
			"tmp/test/hello.txt": &fstest.MapFile{
				Data: []byte("hello"),
				Mode: 0o644,
				Sys:  defaultStat(0o644, 6),
			},
			"tmp/test/world.txt": &fstest.MapFile{
				Data: []byte("world"),
				Mode: 0o644,
				Sys:  defaultStat(0o644, 6),
			},
		}),
		clock.NewFakeClock(time.Date(2000, 12, 15, 17, 8, 0o0, 0, time.UTC)),
		"tmp/test",
	)
	require.NoError(t, err)
	require.NotNil(t, repo)

	initCmd, err := command.NewInitCommand(repo)
	require.NoError(t, err)

	_, err = initCmd.Run()
	require.NoError(t, err)

	statCmd, err := command.NewStatusCommand(repo)
	require.NoError(t, err)

	addCmd, err := command.NewAddCommand([]string{"hello.txt"}, repo)
	require.NoError(t, err)

	_, err = addCmd.Run()
	require.NoError(t, err)

	output, err := statCmd.Run()
	require.NoError(t, err)

	require.EqualValues(t, "?? world.txt\n", string(output))
}

func TestListingUntrackedDirectories(t *testing.T) {
	t.Parallel()

	repo, err := repository.New(
		memory.New(fstest.MapFS{
			"tmp/test/": &fstest.MapFile{
				Mode: os.ModeDir,
			},
			"tmp/test/hello.txt": &fstest.MapFile{
				Data: []byte("hello"),
				Mode: 0o644,
				Sys:  defaultStat(0o644, 6),
			},
			"tmp/test/internal": &fstest.MapFile{
				Mode: os.ModeDir,
			},
			"tmp/test/internal/hello.txt": &fstest.MapFile{
				Data: []byte("hello"),
				Mode: 0o644,
				Sys:  defaultStat(0o644, 6),
			},
			"tmp/test/internal/world.txt": &fstest.MapFile{
				Data: []byte("hello"),
				Mode: 0o644,
				Sys:  defaultStat(0o644, 6),
			},
			"tmp/test/internal/help": &fstest.MapFile{
				Mode: os.ModeDir,
			},
			"tmp/test/internal/help/hello.txt": &fstest.MapFile{
				Data: []byte("hello"),
				Mode: 0o644,
				Sys:  defaultStat(0o644, 6),
			},
		}),
		clock.NewFakeClock(time.Date(2000, 12, 15, 17, 8, 0o0, 0, time.UTC)),
		"tmp/test",
	)
	require.NoError(t, err)
	require.NotNil(t, repo)

	initCmd, err := command.NewInitCommand(repo)
	require.NoError(t, err)

	_, err = initCmd.Run()
	require.NoError(t, err)

	statCmd, err := command.NewStatusCommand(repo)
	require.NoError(t, err)

	addCmd, err := command.NewAddCommand([]string{"hello.txt"}, repo)
	require.NoError(t, err)

	_, err = addCmd.Run()
	require.NoError(t, err)

	output, err := statCmd.Run()
	require.NoError(t, err)

	require.EqualValues(t, "?? internal/\n", string(output))
}
