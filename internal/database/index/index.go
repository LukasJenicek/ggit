package index

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/LukasJenicek/ggit/internal/ds"
	"os"

	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/hasher"
)

const (
	regularMode    = 0o100644
	executableMode = 0o100755
	maxPathSize    = 0xfff
)

type Indexer struct {
	fs         filesystem.Fs
	fileWriter *filesystem.AtomicFileWriter
	database   *database.Database

	indexFilePath string
	rootDir       string
}

func NewIndexer(
	fs filesystem.Fs,
	fileWriter *filesystem.AtomicFileWriter,
	database *database.Database,
	gitPath string,
	rootDir string,
) *Indexer {
	return &Indexer{
		fs:         fs,
		fileWriter: fileWriter,
		database:   database,

		indexFilePath: gitPath + "/index",
		rootDir:       rootDir,
	}
}

func (idx *Indexer) Add(files ds.Set[string]) error {
	if err := idx.createIndexFile(); err != nil {
		return fmt.Errorf("create index file: %w", err)
	}

	indexContent, err := idx.indexContent(files)
	if err != nil {
		return fmt.Errorf("index content: %w", err)
	}

	if err = idx.fileWriter.Update(idx.indexFilePath, indexContent); err != nil {
		return fmt.Errorf("update index file: %w", err)
	}

	return nil
}

func (idx *Indexer) indexContent(files ds.Set[string]) ([]byte, error) {
	indexContent := bytes.NewBuffer(nil)

	entriesLen := files.Size()
	if err := idx.writeHeader(indexContent, entriesLen); err != nil {
		return nil, fmt.Errorf("add: %w", err)
	}

	entries, err := idx.database.SaveBlobs(files)
	if err != nil {
		return nil, fmt.Errorf("save blobs: %w", err)
	}

	indexEntries := make([]*Entry, len(entries))

	for i, e := range entries {
		fInfo, err := idx.fs.Stat(e.Filepath)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", e.Filepath, err)
		}

		indexEntry, err := NewEntry(e.GetRelativeFilePath(idx.rootDir), fInfo, e.OID)
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
func (idx *Indexer) writeHeader(indexContent *bytes.Buffer, entriesLen int) error {
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

// create .git/index file if does not exist.
func (idx *Indexer) createIndexFile() error {
	_, err := idx.fs.Stat(idx.indexFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			if _, err := idx.fs.Create(idx.indexFilePath); err != nil {
				return fmt.Errorf("create index file: %w", err)
			}
		} else {
			return fmt.Errorf("stat %q: %w", idx.indexFilePath, err)
		}
	}

	return nil
}
