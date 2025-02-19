package database

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/LukasJenicek/ggit/internal/filesystem"
)

type Refs struct {
	fs         filesystem.Fs
	fileWriter *filesystem.AtomicFileWriter
	gitDir     string
}

func NewRefs(fs filesystem.Fs, gitDir string, fileWriter *filesystem.AtomicFileWriter) (*Refs, error) {
	if gitDir == "" {
		return nil, errors.New("gitDir is empty")
	}

	if fs == nil {
		return nil, errors.New("fs is required")
	}

	return &Refs{
		fs:         fs,
		gitDir:     gitDir,
		fileWriter: fileWriter,
	}, nil
}

//nolint:wrapcheck
func (r *Refs) UpdateHead(commitID string) error {
	return r.fileWriter.Update(r.headPath(), commitID)
}

func (r *Refs) Current() (string, error) {
	open, err := r.fs.Open(r.headPath())
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return "", fmt.Errorf("could not open HEAD file: %w", err)
		}
	}

	currentCID, err := io.ReadAll(open)
	if err != nil {
		return "", fmt.Errorf("could not read from HEAD file: %w", err)
	}

	return strings.TrimSpace(string(currentCID)), nil
}

func (r *Refs) headPath() string {
	return filepath.Join(r.gitDir, "HEAD")
}
