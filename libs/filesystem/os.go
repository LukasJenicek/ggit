package filesystem

import "os"

type OsFS struct{}

func New() *OsFS {
	return &OsFS{}
}

func (*OsFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (*OsFS) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}
