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
		return fmt.Errorf("opening file %q: %w", path, err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("writing to file %q: %w", lockPath, err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync to file %q: %w", lockPath, err)
	}

	if err = os.Rename(lockPath, path); err != nil {
		return fmt.Errorf("rename tmp object %s: %w", lockPath, err)
	}

	return nil
}
