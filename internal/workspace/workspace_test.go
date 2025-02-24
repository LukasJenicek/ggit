package workspace_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/helpers"
	"github.com/LukasJenicek/ggit/internal/workspace"
)

func TestWorkspace_ListFiles(t *testing.T) {
	t.Parallel()

	projectRootFolder, err := helpers.GetProjectRootFolder()
	if err != nil {
		t.Error(err)
	}

	testDataFolder := filepath.Join(projectRootFolder, "testdata")

	w := workspace.New(testDataFolder, filesystem.New())

	files, err := w.ListFiles(".")
	if err != nil {
		t.Errorf("list files: %s", err)
	}

	expectedFiles := []*workspace.File{
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/a/a.txt",
		},
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/a.txt",
		},
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/b/b.txt",
		},
	}

	assert.EqualValues(t, expectedFiles, files)
}

func TestWorkspace_ListSpecificFiles(t *testing.T) {
	t.Parallel()

	projectRootFolder, err := helpers.GetProjectRootFolder()
	if err != nil {
		t.Error(err)
	}

	testDataFolder := filepath.Join(projectRootFolder, "testdata")

	w := workspace.New(testDataFolder, filesystem.New())

	files, err := w.ListFiles("*.txt")
	if err != nil {
		t.Errorf("list files: %s", err)
	}

	expectedFiles := []*workspace.File{
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/a/a.txt",
		},
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/a.txt",
		},
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/b/b.txt",
		},
	}

	assert.EqualValues(t, expectedFiles, files)
}
