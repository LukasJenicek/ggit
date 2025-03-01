package index

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/hasher"
)

type Content struct {
	database *database.Database
	fs       filesystem.Fs

	rootDir string
}

func (c *Content) Generate(entries []*Entry) ([]byte, error) {
	content := bytes.NewBuffer(nil)

	if err := c.writeHeader(content, len(entries)); err != nil {
		return nil, fmt.Errorf("add: %w", err)
	}

	for _, entry := range entries {
		entryContent, err := entry.Content()
		if err != nil {
			return nil, fmt.Errorf("index entry content: %w", err)
		}

		if err = binary.Write(content, binary.BigEndian, entryContent); err != nil {
			return nil, fmt.Errorf("write index entry: %w", err)
		}
	}

	oid, err := hasher.SHA1HashContent(content.Bytes())
	if err != nil {
		return nil, fmt.Errorf("hash index: %w", err)
	}

	if err = binary.Write(content, binary.BigEndian, oid); err != nil {
		return nil, fmt.Errorf("write index entry: %w", err)
	}

	return content.Bytes(), nil
}

// 12-byte header.
func (c *Content) writeHeader(indexContent *bytes.Buffer, entriesLen int) error {
	header := []any{
		[4]byte{'D', 'I', 'R', 'C'}, // stands for dir cache
		[4]byte{0, 0, 0, 2},         // version
	}

	entries := make([]byte, 4)
	//nolint:gosec
	binary.BigEndian.PutUint32(entries, uint32(entriesLen))
	header = append(header, entries)

	for _, h := range header {
		if err := binary.Write(indexContent, binary.BigEndian, h); err != nil {
			return fmt.Errorf("write header: %w", err)
		}
	}

	return nil
}
