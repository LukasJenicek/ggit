package database

import (
	"crypto/sha1"
	"fmt"
	"os"
	"strings"
)

type Blob struct {
	filename string
	content  []byte
}

func NewBlob(f *os.File, rootDir string) (*Blob, error) {
	content, err := os.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("read file content: %w", err)
	}

	// relative path to root folder
	fPath := strings.Replace(f.Name(), rootDir+"/", "", 1)

	return &Blob{
		filename: fPath,
		content:  content,
	}, nil
}

func (b *Blob) ID() []byte {
	hasher := sha1.New()
	hasher.Write(b.Content())

	return hasher.Sum(nil)
}

func (b *Blob) Type() string {
	return "blob"
}

func (b *Blob) Content() []byte {
	return []byte(fmt.Sprintf("%s %d\x00%s", b.Type(), len(b.content), b.content))
}
