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
	Path string
	File File
}

type FileLocker struct {
	fs Fs
}

func NewFileLocker(fs Fs) *FileLocker {
	return &FileLocker{fs: fs}
}

func (f *FileLocker) Lock(path string) (*LockFile, error) {
	lockPath, err := lockFilePath(path)
	if err != nil {
		return nil, fmt.Errorf("lock File Path: %w", err)
	}

	info, err := f.fs.Stat(lockPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("get lock file info: %w", err)
		}
	}

	if info != nil {
		file, err := f.fs.OpenFile(lockPath, os.O_RDWR|os.O_EXCL, 0o644)
		if err != nil {
			return nil, fmt.Errorf("open lock file: %w", err)
		}

		return &LockFile{File: file, Path: lockPath}, ErrLockAcquired
	}

	lockFile, err := f.fs.OpenFile(lockPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return nil, fmt.Errorf("create lock file %q: %w", lockPath, err)
	}

	return &LockFile{File: lockFile, Path: lockPath}, nil
}

func (f *FileLocker) Unlock(lock *LockFile) error {
	fmt.Println("trying to remove", lock.Path)

	if err := f.fs.Remove(lock.Path); err != nil {
		return fmt.Errorf("remove lock file %q: %w", lock.Path, err)
	}

	return nil
}

func lockFilePath(path string) (string, error) {
	cleanPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve absolute Path: %w", err)
	}

	return cleanPath + ".lock", nil
}
