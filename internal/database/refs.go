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

	headFilePath string
}

func NewRefs(
	fs filesystem.Fs,
	gitDir string,
	fileWriter *filesystem.AtomicFileWriter,
) (*Refs, error) {
	if gitDir == "" {
		return nil, errors.New("gitDir is empty")
	}

	if fs == nil {
		return nil, errors.New("fs is required")
	}

	if fileWriter == nil {
		return nil, errors.New("fileWriter is required")
	}

	return &Refs{
		fs:           fs,
		headFilePath: filepath.Join(gitDir, "HEAD"),
		fileWriter:   fileWriter,
	}, nil
}

func (r *Refs) UpdateHead(commitID string) error {
	if commitID == "" {
		return errors.New("commit id is empty")
	}

	if len(commitID) != 40 {
		return fmt.Errorf("invalid commit ID length: %d", len(commitID))
	}

	if err := r.fileWriter.Update(r.headFilePath, []byte(commitID)); err != nil {
		return fmt.Errorf("update : %w", err)
	}

	return nil
}

func (r *Refs) Current() (string, error) {
	open, err := r.fs.Open(r.headFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return "", fmt.Errorf("open HEAD file: %w", err)
		}
	}

	currentCID, err := io.ReadAll(open)
	if err != nil {
		return "", fmt.Errorf("read from HEAD file: %w", err)
	}

	return strings.TrimSpace(string(currentCID)), nil
}
