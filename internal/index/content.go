package index

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/ds"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/hasher"
)

type Content struct {
	database *database.Database
	fs       filesystem.Fs

	rootDir string
}

func (c *Content) Generate(files ds.Set[string]) ([]byte, error) {
	indexContent := bytes.NewBuffer(nil)

	entriesLen := files.Size()
	if err := c.writeHeader(indexContent, entriesLen); err != nil {
		return nil, fmt.Errorf("add: %w", err)
	}

	entries, err := c.database.SaveBlobs(files)
	if err != nil {
		return nil, fmt.Errorf("save blobs: %w", err)
	}

	indexEntries := make([]*Entry, len(entries))

	for i, e := range entries {
		fInfo, err := c.fs.Stat(e.Filepath)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", e.Filepath, err)
		}

		indexEntry, err := NewEntry(e.GetRelativeFilePath(c.rootDir), fInfo, e.OID)
		if err != nil {
			return nil, fmt.Errorf("new index entry: %w", err)
		}

		indexEntries[i] = indexEntry
	}

	for _, indexEntry := range indexEntries {
		content, err := indexEntry.Content()
		if err != nil {
			return nil, fmt.Errorf("index entry content: %w", err)
		}

		if err = binary.Write(indexContent, binary.BigEndian, content); err != nil {
			return nil, fmt.Errorf("write index entry: %w", err)
		}
	}

	oid, err := hasher.SHA1HashContent(indexContent.Bytes())
	if err != nil {
		return nil, fmt.Errorf("hash index: %w", err)
	}

	if err = binary.Write(indexContent, binary.BigEndian, oid); err != nil {
		return nil, fmt.Errorf("write index entry: %w", err)
	}

	return indexContent.Bytes(), nil
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
