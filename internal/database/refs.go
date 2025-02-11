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
	gitDir string
}

func NewRefs(gitDir string) (*Refs, error) {
	if gitDir == "" {
		return nil, errors.New("gitDir is empty")
	}

	return &Refs{gitDir: gitDir}, nil
}

func (r *Refs) UpdateHead(commitId string) error {
	f, err := os.Create(r.headPath())
	if err != nil {
		return fmt.Errorf("could not open HEAD file: %w", err)
	}
	defer f.Close()

	if _, err = f.WriteString(commitId); err != nil {
		return fmt.Errorf("could not write to HEAD file: %w", err)
	}

	return nil
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
