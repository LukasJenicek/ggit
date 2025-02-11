package database

import (
	"fmt"
	"os"
)

type AtomicFileWriter struct{}

func (*AtomicFileWriter) Update(path string, content string) error {
	lockPath := path + ".lock"

	f, err := os.OpenFile(lockPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("creating temp file %q: %w", lockPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("writing to temp file %q: %w", lockPath, err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync to temp file %q: %w", lockPath, err)
	}

	if err = os.Rename(lockPath, path); err != nil {
		return fmt.Errorf("rename temp file %s: %w", lockPath, err)
	}

	return nil
}
