package command_test

import (
	"bytes"
	"context"
	"io"
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

func TestInitNewRepository(t *testing.T) {
	t.Parallel()

	repo, err := repository.New(
		memory.New(fstest.MapFS{}),
		clock.NewFakeClock(time.Date(2000, 12, 15, 17, 8, 0o0, 0, time.UTC)),
		"tmp/test/",
	)
	require.NoError(t, err)
	require.NotNil(t, repo)

	runner := command.NewRunner(repo)

	osExit, err := runner.RunCmd(context.Background(), "init", []string{}, io.Discard)

	require.NoError(t, err)
	require.Equal(t, 0, osExit)

	g, err := repo.FS.Stat("tmp/test/.git")
	require.NoError(t, err)
	require.True(t, g.IsDir())

	g, err = repo.FS.Stat("tmp/test/.git/refs")
	require.NoError(t, err)
	require.True(t, g.IsDir())

	g, err = repo.FS.Stat("tmp/test/.git/objects")
	require.NoError(t, err)
	require.True(t, g.IsDir())
}

func TestInitAlreadyInitializedRepository(t *testing.T) {
	t.Parallel()

	repo, err := repository.New(
		memory.New(fstest.MapFS{
			"tmp/test/.git": &fstest.MapFile{
				Data:    []byte{},
				Mode:    os.ModeDir,
				ModTime: time.Time{},
				Sys:     nil,
			},
		}),
		clock.NewFakeClock(time.Date(2000, 12, 15, 17, 8, 0o0, 0, time.UTC)),
		"tmp/test/",
	)
	require.NoError(t, err)

	runner := command.NewRunner(repo)

	buf := bytes.NewBuffer(nil)
	_, _ = runner.RunCmd(context.Background(), "init", []string{}, buf)

	require.Equal(t, []byte("Reinitialized existing Git repository in tmp/test/.git"), buf.Bytes())
}
