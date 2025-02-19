package filesystem

import (
	"fmt"
	"os"
)

type AtomicFileWriter struct {
	fs Fs
}

func NewAtomicFileWriter(fs Fs) *AtomicFileWriter {
	return &AtomicFileWriter{fs: fs}
}

func (a *AtomicFileWriter) Update(path string, content string) error {
	lockPath := path + ".lock"

	f, err := a.fs.OpenFile(lockPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create lock file %q: %w", lockPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("write to lock file %q: %w", lockPath, err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync lock file %q: %w", lockPath, err)
	}

	if err = a.fs.Rename(lockPath, path); err != nil {
		return fmt.Errorf("rename lock file %s => %s: %w", lockPath, path, err)
	}

	return nil
}
