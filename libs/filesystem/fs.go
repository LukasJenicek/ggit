package filesystem

import (
	"os"
)

type Dir interface {
	Mkdir(name string, perm os.FileMode) error
}

type Fs interface {
	Stat(name string) (os.FileInfo, error)
	Dir
}
