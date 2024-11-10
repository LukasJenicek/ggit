package memory

import "os"

type MemoryFS struct{}

func New() *MemoryFS {
	return &MemoryFS{}
}

// TODO: Implement
func (*MemoryFS) Stat(name string) (os.FileInfo, error) {
	return nil, nil
}

// TODO: Implement
func (*MemoryFS) Mkdir(name string, perm os.FileMode) error {
	return nil
}
