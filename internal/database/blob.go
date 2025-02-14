package database

import (
	"fmt"
)

type Blob struct {
	filename     string
	content      []byte
	isExecutable bool
}

func NewBlob(content []byte) (*Blob, error) {
	return &Blob{content: content}, nil
}

func (b *Blob) Content() []byte {
	return []byte(fmt.Sprintf("%s %d\x00%s", "blob", len(b.content), b.content))
}
