package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
)

type OsFS struct{}

func New() *OsFS {
	return &OsFS{}
}

//nolint:wrapcheck
func (*OsFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

//nolint:wrapcheck
func (*OsFS) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

//nolint:wrapcheck
func (*OsFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

//nolint:wrapcheck
func (*OsFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(name, flag, perm)
}

//nolint:wrapcheck
func (*OsFS) WalkDir(path string, walkDir fs.WalkDirFunc) error {
	return filepath.WalkDir(path, walkDir)
}

//nolint:wrapcheck
func (*OsFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

//nolint:wrapcheck
func (*OsFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

//nolint:wrapcheck
func (*OsFS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

//nolint:wrapcheck
func (*OsFS) Remove(name string) error {
	return os.Remove(name)
}
