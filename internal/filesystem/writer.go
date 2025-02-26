package filesystem

import (
	"errors"
	"fmt"
	"os"
)

type AtomicFileWriter struct {
	locker Locker
	fs     Fs
}

func NewAtomicFileWriter(fs Fs, locker Locker) (*AtomicFileWriter, error) {
	if fs == nil {
		return nil, errors.New("fs must not be nil")
	}

	if locker == nil {
		return nil, errors.New("locker must not be nil")
	}

	return &AtomicFileWriter{
		fs:     fs,
		locker: locker,
	}, nil
}

func (a *AtomicFileWriter) Write(path string, content []byte) error {
	tmpFilePath := path + ".tmp"

	lock, err := a.locker.Lock(path)
	if err != nil {
		return fmt.Errorf("lock index: %w", err)
	}

	defer func() {
		err = a.locker.Unlock(lock)
	}()

	f, err := a.fs.OpenFile(tmpFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("open tmp file %q: %w", tmpFilePath, err)
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		return fmt.Errorf("write to tmp file %q: %w", tmpFilePath, err)
	}

	if err = f.Sync(); err != nil {
		return fmt.Errorf("sync tmp file %q: %w", tmpFilePath, err)
	}

	if err = a.fs.Rename(tmpFilePath, path); err != nil {
		return fmt.Errorf("rename lock file %s => %s: %w", tmpFilePath, path, err)
	}

	return nil
}
