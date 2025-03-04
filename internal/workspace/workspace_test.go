package workspace_test

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/filesystem/memory"
	"github.com/LukasJenicek/ggit/internal/workspace"
)

func TestWorkspace(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		"tmp/testdata/": &fstest.MapFile{
			Mode: os.ModeDir,
		},
		"tmp/testdata/a.txt": &fstest.MapFile{
			Data: []byte("a"),
			Mode: 0o644,
		},
		"tmp/testdata/b.txt": &fstest.MapFile{
			Data: []byte("b"),
			Mode: 0o644,
		},
		"tmp/testdata/internal/": &fstest.MapFile{
			Mode: os.ModeDir,
		},
		"tmp/testdata/internal/a.txt": &fstest.MapFile{
			Data: []byte("a"),
		},
		"tmp/testdata/internal/b.txt": &fstest.MapFile{
			Data: []byte("b"),
		},
		"tmp/testdata/internal/memory/": &fstest.MapFile{
			Mode: os.ModeDir,
		},
		"tmp/testdata/internal/memory/memory.go": &fstest.MapFile{
			Data: []byte("func main(){ os.Exit(0) }"),
		},
	}

	w, err := workspace.New("tmp/testdata", memory.New(fsys))
	require.NoError(t, err)

	tests := []struct {
		name      string
		pattern   string
		files     []string
		expectErr bool
		err       error
	}{
		{
			name:    "match specific file in current folder",
			pattern: "a.txt",
			files:   []string{"a.txt"},
		},
		{
			name:    "match specific folder and subfolders",
			pattern: "internal/",
			files: []string{
				"internal/a.txt",
				"internal/b.txt",
				"internal/memory/memory.go",
			},
		},
		{
			name:    "match all files",
			pattern: ".",
			files: []string{
				"a.txt",
				"b.txt",
				"internal/a.txt",
				"internal/b.txt",
				"internal/memory/memory.go",
			},
		},
		{
			name:    "match files within current directory with txt extension",
			pattern: "*.txt",
			files: []string{
				"a.txt",
				"b.txt",
			},
		},
		{
			name:      "match zero files should return err",
			pattern:   "blabla.txt",
			files:     []string{},
			expectErr: true,
			err:       &workspace.ErrPathNotMatched{Pattern: "blabla.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := w.ListFiles(tt.pattern)

			if tt.expectErr {
				require.Error(t, err)
				require.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, tt.files, f)
			}
		})
	}
}
