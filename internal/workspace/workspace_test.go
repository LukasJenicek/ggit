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

	assert.EqualValues(
		t,
		[]string{
			filepath.Join(testDataFolder, "a", "a.txt"),
			filepath.Join(testDataFolder, "a.txt"),
			filepath.Join(testDataFolder, "b", "b.txt"),
		}, files,
	)
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

	assert.EqualValues(
		t,
		[]string{
			filepath.Join(testDataFolder, "a", "a.txt"),
			filepath.Join(testDataFolder, "a.txt"),
			filepath.Join(testDataFolder, "b", "b.txt"),
		}, files,
	)
}

func TestWorkspace_ListSpecificFile(t *testing.T) {
	t.Parallel()

	projectRootFolder, err := helpers.GetProjectRootFolder()
	if err != nil {
		t.Error(err)
	}

	testDataFolder := filepath.Join(projectRootFolder, "testdata")

	w := workspace.New(testDataFolder, filesystem.New())

	files, err := w.ListFiles("a.txt")
	if err != nil {
		t.Errorf("list files: %s", err)
	}

	assert.EqualValues(
		t,
		[]string{
			filepath.Join(testDataFolder, "a", "a.txt"),
			filepath.Join(testDataFolder, "a.txt"),
		}, files,
	)
}
