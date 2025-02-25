package filesystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrLockAcquired = errors.New("lock is acquired")

type Locker interface {
	Lock(path string) (*LockFile, error)
	Unlock(lock *LockFile) error
}

type LockFile struct {
	path string
	file File
}

type FileLocker struct {
	fs Fs
}

func NewFileLocker(fs Fs) *FileLocker {
	return &FileLocker{fs: fs}
}

func (f *FileLocker) Lock(path string) (*LockFile, error) {
	lockFilePath, err := lockFilePath(path)
	if err != nil {
		return nil, fmt.Errorf("lock file path: %w", err)
	}

	info, err := f.fs.Stat(lockFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("get lock file info: %w", err)
		}
	}

	if info != nil {
		file, err := f.fs.OpenFile(lockFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return nil, fmt.Errorf("open lock file: %w", err)
		}

		return &LockFile{file: file, path: lockFilePath}, ErrLockAcquired
	}

	lockFile, err := f.fs.OpenFile(lockFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return nil, fmt.Errorf("create lock file %q: %w", lockFilePath, err)
	}

	return &LockFile{file: lockFile, path: lockFilePath}, nil
}

func (f *FileLocker) Unlock(lock *LockFile) error {
	defer lock.file.Close()

	lockFilePath, err := lockFilePath(lock.path)
	if err != nil {
		return fmt.Errorf("lock file path: %w", err)
	}

	if err = f.fs.Remove(lockFilePath); err != nil {
		return fmt.Errorf("remove lock file %q: %w", lockFilePath, err)
	}

	return nil
}

func lockFilePath(path string) (string, error) {
	cleanPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path: %w", err)
	}

	return cleanPath + ".lock", nil
}
