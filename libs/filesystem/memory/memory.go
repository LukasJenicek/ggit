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

// GetFs
// Testing purposes only
func (f *Fs) Get() map[string]os.FileInfo {
	return f.fs
}

// Stat simulates the os.Stat function for the in-memory file system.
// It returns the file info if the file or directory exists, or an error if it does not.
func (f *Fs) Stat(name string) (os.FileInfo, error) {
	info, exists := f.fs[name]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", name)
	}
	return info, nil
}

// Mkdir simulates the os.Mkdir function for the in-memory file system.
// It creates a new directory in the in-memory file system.
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

// CreateFile creates a new file in the in-memory file system with initial content.
func (f *Fs) CreateFile(name string, content []byte) error {
	if _, exists := f.fs[name]; exists {
		return fmt.Errorf("file already exists: %s", name)
	}

	// Create a new file (represented by a FileInfo object)
	f.fs[name] = &FileInfo{
		name:    name,
		isDir:   false,
		content: content,
	}

	// Store the file content in the memory
	f.contents[name] = content

	return nil
}

// WriteFile writes content to an existing file in the in-memory file system.
func (f *Fs) WriteFile(name string, content []byte) error {
	if _, exists := f.fs[name]; !exists {
		return fmt.Errorf("file not found: %s", name)
	}

	// Update the content in memory
	f.contents[name] = content
	return nil
}

// ReadFile reads the content of a file in the in-memory file system.
func (f *Fs) ReadFile(name string) ([]byte, error) {
	content, exists := f.contents[name]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", name)
	}
	return content, nil
}

// FileInfo is a custom implementation of os.FileInfo for in-memory file system.
type FileInfo struct {
	name    string
	isDir   bool
	content []byte // file content for non-directory files
}

// Name returns the name of the file or directory.
func (fi *FileInfo) Name() string {
	return fi.name
}

// Size returns the size of the file. For directories, we return 0.
func (fi *FileInfo) Size() int64 {
	if fi.isDir {
		return 0
	}
	return int64(len(fi.content)) // Return the content size for files
}

// Mode returns the file mode (permissions). For directories, we return os.ModeDir.
func (fi *FileInfo) Mode() os.FileMode {
	if fi.isDir {
		return os.ModeDir
	}
	return 0 // Default, can be expanded
}

// ModTime returns the last modification time. For simplicity, we use the current time.
func (fi *FileInfo) ModTime() time.Time {
	return time.Now()
}

// IsDir returns true if the file is a directory, false otherwise.
func (fi *FileInfo) IsDir() bool {
	return fi.isDir
}

// Sys returns nil (can be expanded if needed for additional attributes).
func (fi *FileInfo) Sys() interface{} {
	return nil
}
