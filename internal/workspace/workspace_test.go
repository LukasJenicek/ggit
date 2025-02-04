package workspace_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LukasJenicek/ggit/internal/workspace"

	"github.com/stretchr/testify/assert"
)

func TestWorkspace_ListFiles(t *testing.T) {
	w := workspace.New()

	projectRootFolder, err := getProjectRootFolder()
	if err != nil {
		t.Error(err)
	}

	testDataFolder := filepath.Join(projectRootFolder, "testdata")

	files, err := w.ListFiles(testDataFolder)
	if err != nil {
		t.Errorf("list files: %s", err)
	}

	expectedFiles := []string{
		testDataFolder + "/a/a.txt",
		testDataFolder + "/a.txt",
		testDataFolder + "/b/b.txt",
	}

	assert.EqualValues(t, expectedFiles, files)
}

func getProjectRootFolder() (string, error) {
	getwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %s", err)
	}

	return strings.Replace(getwd, "/internal/workspace", "", 1), nil
}
