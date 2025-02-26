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
	gitDir       string
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
		gitDir:       gitDir,
	}, nil
}

func (r *Refs) InitRef(ref string) error {
	if ref == "" {
		return errors.New("ref is empty")
	}

	if err := r.fileWriter.Write(r.headFilePath, []byte(ref)); err != nil {
		return fmt.Errorf("update : %w", err)
	}

	return nil
}

func (r *Refs) UpdateHead(commitID string) error {
	if commitID == "" {
		return errors.New("commit id is empty")
	}

	currentRef, err := r.parseCurrentRef()
	if err != nil {
		return fmt.Errorf("parse current ref: %w", err)
	}

	refPath := filepath.Join(r.gitDir, "refs", "heads", currentRef)

	if err := r.fileWriter.Write(refPath, []byte(commitID)); err != nil {
		return fmt.Errorf("update : %w", err)
	}

	return nil
}

func (r *Refs) Current() (string, error) {
	currentRef, err := r.parseCurrentRef()
	if err != nil {
		return "", fmt.Errorf("parse current ref: %w", err)
	}

	refPath := filepath.Join(r.gitDir, "refs", "heads", currentRef)

	refFile, err := r.fs.Open(refPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return "", fmt.Errorf("open HEAD file: %w", err)
		}
	}

	currentCID, err := io.ReadAll(refFile)
	if err != nil {
		return "", fmt.Errorf("read from HEAD file: %w", err)
	}

	return strings.TrimSpace(string(currentCID)), nil
}

func (r *Refs) parseCurrentRef() (string, error) {
	refs, err := r.fs.ReadFile(r.headFilePath)
	if err != nil {
		return "", fmt.Errorf("read refs: %w", err)
	}

	return strings.Replace(string(refs), "ref: refs/heads/", "", 1), nil
}
