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

type Index struct {
	Entries *Entries
	Parents *Parents
}

func NewIndex() *Index {
	return &Index{
		Entries: NewEntries(),
		Parents: NewParents(),
	}
}

func (i *Index) Tracked(path string) bool {
	if _, ok := i.Entries.Get(path); ok {
		return true
	}

	if _, ok := i.Parents.Get(path); ok {
		return true
	}

	return false
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

// Add
// Start tracking files using .git/index.
func (i *Indexer) Add(files []string) error {
	index, err := i.Load()
	if err != nil {
		return fmt.Errorf("load index: %w", err)
	}

	entries, err := i.database.SaveBlobs(ds.NewSet(files))
	if err != nil {
		return fmt.Errorf("save blobs: %w", err)
	}

	if index.Entries.Len() > 0 {
		i.clean(files, index.Entries, index.Parents)
	}

	for _, e := range entries {
		stat, err := i.fs.Stat(e.AbsFilePath)
		if err != nil {
			return fmt.Errorf("stat %s: %w", e.AbsFilePath, err)
		}

		relFilePath := e.GetRelativeFilePath(i.rootDir)

		indexEntry, err := NewEntry(relFilePath, stat, e.OID)
		if err != nil {
			return fmt.Errorf("new index entry: %w", err)
		}

		index.Entries.Add(relFilePath, indexEntry)
	}

	indexContent, err := i.content.Generate(index.Entries.SortedValues())
	if err != nil {
		return fmt.Errorf("index content: %w", err)
	}

	if err := i.createIndexFile(); err != nil {
		return fmt.Errorf("create index file: %w", err)
	}

	if err = i.fileWriter.Write(i.indexFilePath, indexContent); err != nil {
		return fmt.Errorf("update index file: %w", err)
	}

	return nil
}

func (i *Indexer) Load() (*Index, error) {
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
		// this is valid case when user is adding files for the first time
		if errors.Is(err, os.ErrNotExist) {
			return NewIndex(), nil
		}

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

	index := NewIndex()

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

		dir, _ := filepath.Split(string(entry.Path))

		parentsFolderPaths := strings.Split(dir, string(os.PathSeparator))
		for i, parentsFolderPath := range parentsFolderPaths {
			if parentsFolderPath == "" {
				continue
			}

			p := parentsFolderPath
			if i > 0 {
				p = strings.Join(parentsFolderPaths[:i+1], string(os.PathSeparator))
			}

			index.Parents.Add(p, entry)
		}

		index.Entries.Add(string(entry.Path), entry)

		cursorPos++
		// finding last null byte of entry
		for cursorPos > 0 && content[cursorPos-1] == 0 {
			cursorPos++
		}

		cursorPos--
		currPosition = cursorPos
	}

	return index, nil
}

// that file must be deleted.
func (i *Indexer) clean(filePaths []string, indexEntries *Entries, parents *Parents) {
	for _, filePath := range filePaths {
		filepathParts := strings.Split(filePath, string(os.PathSeparator))

		for i, fPart := range filepathParts {
			part := fPart
			if i > 0 {
				part = filepath.Join(filepathParts[:i]...)
			}

			indexEntries.Delete(part)

			// if new file conflicts with existing folder
			// all files within that folder must be deleted
			if entries, ok := parents.Get(part); ok {
				for _, entry := range entries {
					indexEntries.Delete(string(entry.Path))
				}

				parents.Delete(part)
			}
		}
	}
}

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
