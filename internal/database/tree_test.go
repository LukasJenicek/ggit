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
	entries = append(entries, database.NewEntry("hello.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false))
	entries = append(entries, database.NewEntry("libs/hello.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false))
	entries = append(entries, database.NewEntry("libs/internal/internal.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false))
	entries = append(entries, database.NewEntry("world.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false))

	want := database.NewTree(nil, "")
	want.AddEntry(database.NewEntry("hello.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false))

	libsTree := database.NewTree(want, "libs")
	libsTree.AddEntry(database.NewEntry("libs/hello.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false))
	libsInternalTree := database.NewTree(libsTree, "internal")
	libsInternalTree.AddEntry(database.NewEntry("libs/internal/internal.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false))
	libsTree.AddEntry(libsInternalTree)
	want.AddEntry(libsTree)
	want.AddEntry(database.NewEntry("world.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false))

	build, err := database.Build(root, entries)

	require.NoError(t, err)
	require.EqualValues(t, want, build)
}
