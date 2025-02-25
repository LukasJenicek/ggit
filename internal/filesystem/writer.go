package filesystem

import (
	"errors"
	"fmt"
)

type AtomicFileWriter struct {
	fs     Fs
	locker Locker
}

func NewAtomicFileWriter(fs Fs) *AtomicFileWriter {
	return &AtomicFileWriter{fs: fs, locker: NewFileLocker(fs)}
}

func (a *AtomicFileWriter) Write(path string, content []byte) error {
	lock, err := a.locker.Lock(path)
	if err != nil {
		if errors.Is(err, ErrLockAcquired) {
			// TODO: either err out or remove lock based on `modification time` ?
			// do something
		} else {
			return fmt.Errorf("lock file %q: %w", path, err)
		}
	}

	defer func() {
		err = a.locker.Unlock(lock)
	}()

	if _, err = lock.file.Write(content); err != nil {
		return fmt.Errorf("write to lock file %q: %w", lock.path, err)
	}

	if err = lock.file.Sync(); err != nil {
		return fmt.Errorf("sync lock file %q: %w", lock.path, err)
	}

	if err = a.fs.Rename(lock.path, path); err != nil {
		return fmt.Errorf("rename lock file %s => %s: %w", lock.path, path, err)
	}

	return nil
}
