package database_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/database"
)

func TestBuildTree(t *testing.T) {
	t.Parallel()

	root := database.NewTree(nil, "")

	entries := []*database.Entry{}

	content := []byte{0xde, 0xad, 0xbe, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xad, 0xde, 0xa, 0xbe, 0xef, 0xad, 0xad, 0xef, 0xef, 0xde}

	entry, err := database.NewEntry("hello.txt", content, false)
	require.NoError(t, err)
	entries = append(entries, entry)

	entry, err = database.NewEntry("libs/hello.txt", content, false)
	require.NoError(t, err)
	entries = append(entries, entry)

	entry, err = database.NewEntry("libs/internal/internal.txt", content, false)
	require.NoError(t, err)
	entries = append(entries, entry)

	entry, err = database.NewEntry("world.txt", content, false)
	require.NoError(t, err)
	entries = append(entries, entry)

	want := database.NewTree(nil, "")
	entry, err = database.NewEntry("hello.txt", content, false)
	require.NoError(t, err)
	want.AddEntry(entry)

	libsTree := database.NewTree(want, "libs")
	entry, err = database.NewEntry("libs/hello.txt", content, false)
	require.NoError(t, err)
	libsTree.AddEntry(entry)

	libsInternalTree := database.NewTree(libsTree, "internal")
	entry, err = database.NewEntry("libs/internal/internal.txt", content, false)
	require.NoError(t, err)
	libsInternalTree.AddEntry(entry)
	libsTree.AddEntry(libsInternalTree)
	want.AddEntry(libsTree)
	entry, err = database.NewEntry("world.txt", content, false)
	require.NoError(t, err)
	want.AddEntry(entry)

	build, err := database.Build(root, entries)

	require.NoError(t, err)
	require.EqualValues(t, want, build)
}
