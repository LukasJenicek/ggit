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

	docsTree := database.NewTree(root, "docs")
	docsTree.AddEntry(
		database.NewEntry("docs.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false),
	)

	root.AddEntry(docsTree)
	root.AddEntry(
		database.NewEntry("hello.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false),
	)

	libsTree := database.NewTree(root, "libs")
	libsTree.AddEntry(
		database.NewEntry("hello.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false),
	)

	libsInternalTree := database.NewTree(libsTree, "libs/internal")
	libsInternalTree.AddEntry(
		database.NewEntry("internal.txt", []byte{0xde, 0xad, 0xbe, 0xef}, false),
	)
	libsTree.AddEntry(libsInternalTree)

	root.AddEntry(libsTree)

	_, err := d.StoreTree(root)

	require.NoError(t, err)
}
