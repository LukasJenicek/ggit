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

func (a *AtomicFileWriter) Write(path string, content []byte) error {
	tmpFilePath := path + ".tmp"

	f, err := a.fs.OpenFile(tmpFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("open tmp file %q: %w", tmpFilePath, err)
	}

	if _, err := f.Write(content); err != nil {
		return fmt.Errorf("write to tmp file %q: %w", tmpFilePath, err)
	}

	if err = f.Sync(); err != nil {
		return fmt.Errorf("sync lock file %q: %w", tmpFilePath, err)
	}

	if err = a.fs.Rename(tmpFilePath, path); err != nil {
		return fmt.Errorf("rename lock file %s => %s: %w", tmpFilePath, path, err)
	}

	return nil
}
