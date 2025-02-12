package workspace_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

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

	w := workspace.New(testDataFolder)

	files, err := w.ListFiles()
	if err != nil {
		t.Errorf("list files: %s", err)
	}

	expectedFiles := []*workspace.File{
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/a",
			Dir:  true,
		},
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/a/a.txt",
			Dir:  false,
		},
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/a.txt",
			Dir:  false,
		},
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/b",
			Dir:  true,
		},
		{
			Path: "/home/lj/Projects/LukasJenicek/ggit/testdata/b/b.txt",
			Dir:  false,
		},
	}

	assert.EqualValues(t, expectedFiles, files)
}
