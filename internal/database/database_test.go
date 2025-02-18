package database_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/database"
)

func TestDatabase_StoreRootTree(t *testing.T) {
	t.Parallel()

	d := database.New("/tmp/.git")

	root := database.NewTree(nil, "")

	content := []byte{0xde, 0xad, 0xbe, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xef, 0xad, 0xde, 0xa, 0xbe, 0xef, 0xad, 0xad, 0xef, 0xef, 0xde}

	docsTree := database.NewTree(root, "docs")
	entry, err := database.NewEntry("docs.txt", content, false)
	require.NoError(t, err)
	docsTree.AddEntry(entry)

	root.AddEntry(docsTree)
	entry, err = database.NewEntry("hello.txt", content, false)
	require.NoError(t, err)
	root.AddEntry(entry)

	libsTree := database.NewTree(root, "libs")
	entry, err = database.NewEntry("hello.txt", content, false)
	require.NoError(t, err)
	libsTree.AddEntry(entry)

	libsInternalTree := database.NewTree(libsTree, "libs/internal")
	entry, err = database.NewEntry("internal.txt", content, false)
	require.NoError(t, err)
	libsInternalTree.AddEntry(entry)
	libsTree.AddEntry(libsInternalTree)

	root.AddEntry(libsTree)

	_, err = d.StoreTree(root)
	require.NoError(t, err)
}
