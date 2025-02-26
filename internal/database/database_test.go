package database_test

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/filesystem/memory"
)

func TestDatabase_StoreRootTree(t *testing.T) {
	t.Parallel()

	d, err := database.New(memory.New(fstest.MapFS{}), "tmp/.git")
	require.NoError(t, err)

	root := database.NewTree(nil, "")

	content := []byte{0xde, 0xad, 0xbe, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xad, 0xde, 0xa, 0xbe, 0xef, 0xad, 0xad, 0xef, 0xef, 0xde}

	docsTree := database.NewTree(root, "docs")
	entry, err := database.NewEntry("docs.txt", "tmp/docs.txt", content, false)
	require.NoError(t, err)
	docsTree.AddEntry(entry)

	root.AddEntry(docsTree)

	entry, err = database.NewEntry("hello.txt", "tmp/hello.txt", content, false)
	require.NoError(t, err)
	root.AddEntry(entry)

	libsTree := database.NewTree(root, "libs")
	entry, err = database.NewEntry("hello.txt", "tmp/libs/hello.txt", content, false)
	require.NoError(t, err)
	libsTree.AddEntry(entry)

	libsInternalTree := database.NewTree(libsTree, "libs/internal")
	entry, err = database.NewEntry("internal.txt", "tmp/libs/internal/internal.txt", content, false)
	require.NoError(t, err)
	libsInternalTree.AddEntry(entry)
	libsTree.AddEntry(libsInternalTree)

	root.AddEntry(libsTree)

	oid, err := d.StoreTree(root)
	require.NoError(t, err)
	require.NotEmpty(t, oid, "object ID should not be empty")
}
