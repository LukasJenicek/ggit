package memory

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing/fstest"

	"github.com/LukasJenicek/ggit/internal/filesystem"
)

type Fs struct {
	fsys fstest.MapFS
}

type MockFile struct {
	content  strings.Builder
	isClosed bool
}

func (m *MockFile) WriteString(s string) (int, error) {
	if m.isClosed {
		return 0, errors.New("write to closed file")
	}

	return m.content.WriteString(s)
}

func (m *MockFile) Sync() error {
	if m.isClosed {
		return errors.New("sync on closed file")
	}

	return nil
}

func (m *MockFile) Close() error {
	m.isClosed = true

	return nil
}

// New creates a new in-memory file system.
func New(fsys fstest.MapFS) *Fs {
	return &Fs{
		fsys: fsys,
	}
}

func (f *Fs) Stat(name string) (os.FileInfo, error) {
	return f.fsys.Stat(name)
}

func (f *Fs) Open(name string) (fs.File, error) {
	return f.fsys.Open(name)
}

func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (filesystem.File, error) {
	_, err := f.fsys.Open(name)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	if errors.Is(err, fs.ErrNotExist) {
		f.fsys[name] = &fstest.MapFile{Data: nil, Mode: perm}
	}

	return &MockFile{
		content:  strings.Builder{},
		isClosed: false,
	}, nil
}

func (f *Fs) WalkDir(path string, walkDir fs.WalkDirFunc) error {
	return fs.WalkDir(f.fsys, path, walkDir)
}

func (f *Fs) ReadFile(name string) ([]byte, error) {
	return f.fsys.ReadFile(name)
}

func (f *Fs) Mkdir(name string, perm os.FileMode) error {
	if _, err := f.fsys.Open(name); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("open: %s", name)
		}
	}

	f.fsys[name] = &fstest.MapFile{Mode: fs.ModeDir}

	return nil
}

func (f *Fs) WriteFile(name string, data []byte, perm os.FileMode) error {
	f.fsys[name] = &fstest.MapFile{Data: data, Mode: fs.ModePerm}

	return nil
}

func (f *Fs) Rename(oldpath, newpath string) error {
	file, ok := f.fsys[oldpath]
	if !ok {
		return os.ErrNotExist
	}

	// errors out when the file with the same name already exist
	_, ok = f.fsys[newpath]
	if ok {
		return os.ErrExist
	}

	delete(f.fsys, oldpath)

	f.fsys[newpath] = file

	return nil
}

func (f *Fs) Remove(name string) error {
	_, ok := f.fsys[name]
	if !ok {
		return os.ErrNotExist
	}

	delete(f.fsys, name)
	return nil
}
