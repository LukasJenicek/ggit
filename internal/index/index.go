package index

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/LukasJenicek/ggit/internal/database"
	"github.com/LukasJenicek/ggit/internal/ds"
	"github.com/LukasJenicek/ggit/internal/filesystem"
)

const (
	regularMode    = 0o100644
	executableMode = 0o100755
	maxPathSize    = 0xfff
)

type Indexer struct {
	fs            filesystem.Fs
	fileWriter    *filesystem.AtomicFileWriter
	database      *database.Database
	locker        filesystem.Locker
	content       *Content
	indexFilePath string
	rootDir       string
}

func NewIndexer(
	fs filesystem.Fs,
	fileWriter *filesystem.AtomicFileWriter,
	locker filesystem.Locker,
	database *database.Database,
	rootDir string,
) (*Indexer, error) {
	if fs == nil {
		return nil, errors.New("file system is nil")
	}

	if fileWriter == nil {
		return nil, errors.New("file writer is nil")
	}

	if locker == nil {
		return nil, errors.New("locker is nil")
	}

	if database == nil {
		return nil, errors.New("database is nil")
	}

	if rootDir == "" {
		return nil, errors.New("root dir is empty")
	}

	return &Indexer{
		fs:         fs,
		fileWriter: fileWriter,
		database:   database,
		locker:     locker,
		content: &Content{
			database: database,
			fs:       fs,
			rootDir:  rootDir,
		},
		rootDir:       rootDir,
		indexFilePath: filepath.Join(rootDir, ".git", "index"),
	}, nil
}

func (i *Indexer) Add(files []string) error {
	if err := i.createIndexFile(); err != nil {
		return fmt.Errorf("create index file: %w", err)
	}

	indexEntries, err := i.LoadIndex()
	if err != nil {
		return fmt.Errorf("load index: %w", err)
	}

	entries, err := i.database.SaveBlobs(ds.NewSet(files))
	if err != nil {
		return fmt.Errorf("save blobs: %w", err)
	}

	i.discardConflicts(files, indexEntries)

	for _, e := range entries {
		fInfo, err := i.fs.Stat(e.Filepath)
		if err != nil {
			return fmt.Errorf("stat %s: %w", e.Filepath, err)
		}

		relativeFilePath := e.GetRelativeFilePath(i.rootDir)

		indexEntry, err := NewEntry(relativeFilePath, fInfo, e.OID)
		if err != nil {
			return fmt.Errorf("new index entry: %w", err)
		}

		indexEntries[relativeFilePath] = indexEntry
	}

	indexContent, err := i.content.Generate(indexEntries.SortedValues())
	if err != nil {
		return fmt.Errorf("index content: %w", err)
	}

	if err = i.fileWriter.Write(i.indexFilePath, indexContent); err != nil {
		return fmt.Errorf("update index file: %w", err)
	}

	return nil
}

// hello.txt/world.txt and hello.txt file can't coexist in the same index.
func (i *Indexer) discardConflicts(files []string, indexEntries Entries) {
	for _, f := range files {
		filepathParts := strings.Split(f, string(os.PathSeparator))
		if len(filepathParts) == 1 {
			continue
		}

		for i, fPart := range filepathParts {
			part := fPart
			if i > 0 {
				part = filepath.Join(filepathParts[:i]...)
			}

			delete(indexEntries, part)
		}
	}
}

// TODO: Determine handling strategy when no files are added to the index.
// Options: return error, create empty index, or maintain current behavior.
func (i *Indexer) LoadIndex() (Entries, error) {
	lock, err := i.locker.Lock(i.indexFilePath)
	if err != nil {
		return nil, fmt.Errorf("lock index: %w", err)
	}

	defer func() {
		err = i.locker.Unlock(lock)
		if err != nil {
			// TODO: log error ??
		}
	}()

	indexFile, err := i.fs.Open(i.indexFilePath)
	if err != nil {
		return nil, fmt.Errorf("open index file %q: %w", i.indexFilePath, err)
	}

	content, err := io.ReadAll(indexFile)
	if err != nil {
		return nil, fmt.Errorf("read index file: %w", err)
	}

	if len(content) > 0 {
		if err = CheckIndexIntegrity(content); err != nil {
			return nil, fmt.Errorf("index integrity: %w", err)
		}
	}

	entryLen := binary.BigEndian.Uint32(content[8:12])
	entries := make(Entries, entryLen)

	currPosition := 12
	for range entryLen {
		pathLen := binary.BigEndian.Uint16(content[currPosition+60 : currPosition+62])
		cursorPos := currPosition + 62

		if pathLen < maxPathSize {
			cursorPos += int(pathLen)
		} else {
			pathLen = 0

			for content[cursorPos] != 0 {
				cursorPos++
				pathLen++
			}
		}

		entry, err := NewEntryFromBytes(content[currPosition:cursorPos], int(pathLen))
		if err != nil {
			return nil, fmt.Errorf("create entry: %w", err)
		}

		entries[string(entry.Path)] = entry
		cursorPos++
		// finding last null byte of entry
		for cursorPos > 0 && content[cursorPos-1] == 0 {
			cursorPos++
		}

		cursorPos--

		currPosition = cursorPos
	}

	return entries, nil
}

// create .git/index file if does not exist.
func (i *Indexer) createIndexFile() error {
	_, err := i.fs.Stat(i.indexFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			if _, err := i.fs.Create(i.indexFilePath); err != nil {
				return fmt.Errorf("create index file: %w", err)
			}
		} else {
			return fmt.Errorf("stat %q: %w", i.indexFilePath, err)
		}
	}

	return nil
}
