package index

import (
	"fmt"
	"github.com/LukasJenicek/ggit/internal/filesystem"
	"os"
)

const regularMode = 0100644
const executableMode = 0100755
const maxPathSize = 0xfff

type Index struct {
	fs         filesystem.Fs
	fileWriter *filesystem.AtomicFileWriter

	filePath string
}

func NewIndex(gitPath string) *Index {
	return &Index{
		filePath: gitPath + "/index",
	}
}

func (idx *Index) Add(path string, oid string) error {
	if err := idx.createIndex(); err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	return nil
}

// createIndex
// create .git/index file if does not exist
func (idx *Index) createIndex() error {
	_, err := idx.fs.Stat(idx.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			if _, err := idx.fs.Create(idx.filePath); err != nil {
				return fmt.Errorf("create index file: %w", err)
			}
		} else {
			return fmt.Errorf("stat %q: %w", idx.filePath, err)
		}
	}

	return nil
}
