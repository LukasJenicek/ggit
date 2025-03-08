package database_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/hasher"
)

func TestBuild(t *testing.T) {
	content, err := hasher.SHA1HashContent([]byte("hello"))
	require.NoError(t, err)

	entries := []*database.Entry{
		{
			Name:        "hello.txt",
			AbsFilePath: "hello.txt",
			OID:         content,
			Executable:  false,
		},
		{
			Name:        "world.txt",
			AbsFilePath: "internal/help/world.txt",
			OID:         content,
			Executable:  false,
		},
		{
			Name:        "world.txt",
			AbsFilePath: "internal/world.txt",
			OID:         content,
			Executable:  false,
		},
	}

	expectedTree := &database.Tree{}

	entry, err := database.NewEntry("hello.txt", "hello.txt", content, false)
	require.NoError(t, err)
	expectedTree.AddEntry(entry)

	internalTree := database.NewTree(expectedTree, "internal")
	expectedTree.AddEntry(internalTree)

	helpTree := database.NewTree(internalTree, "help")
	entry, err = database.NewEntry("world.txt", "internal/help/world.txt", content, false)
	require.NoError(t, err)
	helpTree.AddEntry(entry)

	internalTree.AddEntry(helpTree)

	entry, err = database.NewEntry("world.txt", "internal/world.txt", content, false)
	require.NoError(t, err)
	internalTree.AddEntry(entry)

	got, err := database.Build(database.NewRootTree(), entries)
	require.NoError(t, err)
	require.EqualValues(t, expectedTree, got)
}
