package filesystem

import (
	"io"
	"io/fs"
	"os"
)

type Dir interface {
	Mkdir(name string, perm os.FileMode) error
	WalkDir(path string, walkDir fs.WalkDirFunc) error
}

type FileWriter interface {
	Create(name string) (io.WriteCloser, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	Rename(oldpath, newpath string) error
	Remove(name string) error
}

type Fs interface {
	fs.FS
	fs.StatFS
	fs.ReadFileFS
	OpenFile(name string, flag int, perm os.FileMode) (File, error)

	FileWriter
	Dir
}

type File interface {
	io.Closer

	Write(b []byte) (int, error)
	Sync() error
}
