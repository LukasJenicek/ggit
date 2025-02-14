package database

import (
	"fmt"
)

type Blob struct {
	content []byte
}

func NewBlob(content []byte) *Blob {
	return &Blob{content: content}
}

func (b *Blob) Content() ([]byte, error) {
	return []byte(fmt.Sprintf("%s %d\x00%s", "blob", len(b.content), b.content)), nil
}
