package memory

import (
	"fmt"
	"os"
	"time"
)

// Fs represents a simple in-memory file system that supports file contents.
type Fs struct {
	fs map[string]os.FileInfo
	// Store file content (in-memory)
	contents map[string][]byte
}

// New creates a new in-memory file system.
func New() *Fs {
	return &Fs{
		fs:       make(map[string]os.FileInfo),
		contents: make(map[string][]byte),
	}
}

func (f *Fs) Get() map[string]os.FileInfo {
	return f.fs
}

func (f *Fs) Stat(name string) (os.FileInfo, error) {
	info, exists := f.fs[name]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", name)
	}

	return info, nil
}

func (f *Fs) Mkdir(name string, perm os.FileMode) error {
	if _, exists := f.fs[name]; exists {
		return fmt.Errorf("directory already exists: %s", name)
	}

	// Create a new directory (represented by a dummy FileInfo object)
	f.fs[name] = &FileInfo{
		name:  name,
		isDir: true,
	}

	return nil
}

func (f *Fs) CreateFile(name string, content []byte) error {
	if _, exists := f.fs[name]; exists {
		return fmt.Errorf("file already exists: %s", name)
	}

	f.fs[name] = &FileInfo{
		name:    name,
		isDir:   false,
		content: content,
	}

	f.contents[name] = content

	return nil
}

func (f *Fs) WriteFile(name string, content []byte) error {
	if _, exists := f.fs[name]; !exists {
		return fmt.Errorf("file not found: %s", name)
	}

	f.contents[name] = content

	return nil
}

func (f *Fs) ReadFile(name string) ([]byte, error) {
	content, exists := f.contents[name]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", name)
	}

	return content, nil
}

type FileInfo struct {
	name    string
	isDir   bool
	content []byte
}

func (fi *FileInfo) Name() string {
	return fi.name
}

func (fi *FileInfo) Size() int64 {
	if fi.isDir {
		return 0
	}

	return int64(len(fi.content))
}

func (fi *FileInfo) Mode() os.FileMode {
	if fi.isDir {
		return os.ModeDir
	}

	return 0
}

func (fi *FileInfo) ModTime() time.Time {
	return time.Now()
}

func (fi *FileInfo) IsDir() bool {
	return fi.isDir
}

func (fi *FileInfo) Sys() interface{} {
	return nil
}
