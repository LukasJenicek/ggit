package database

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Refs struct {
	fileWriter *AtomicFileWriter

	gitDir string
}

func NewRefs(gitDir string, fileWriter *AtomicFileWriter) (*Refs, error) {
	if gitDir == "" {
		return nil, errors.New("gitDir is empty")
	}

	return &Refs{gitDir: gitDir, fileWriter: fileWriter}, nil
}

func (r *Refs) UpdateHead(commitId string) error {
	return r.fileWriter.Update(r.headPath(), commitId)
}

func (r *Refs) Current() (string, error) {
	open, err := os.Open(r.headPath())
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
