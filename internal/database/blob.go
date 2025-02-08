package database

import (
	"crypto/sha1"
	"fmt"
	"os"
)

type Blob struct {
	filename string
	content  []byte
}

func NewBlob(filePath string) (*Blob, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", filePath, err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file content: %w", err)
	}

	return &Blob{
		filename: f.Name(),
		content:  content,
	}, nil
}

func (b *Blob) Id() []byte {
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
