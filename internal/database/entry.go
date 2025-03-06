package database

import (
	"errors"
	"fmt"
	"strings"
)

type Entry struct {
	// Name of the file
	Name string
	// Absolute filepath
	AbsFilePath string
	// Calculated object id
	OID        []byte
	Executable bool
}

func NewEntry(filename string, absFilePath string, oid []byte, executable bool) (*Entry, error) {
	if strings.TrimSpace(filename) == "" {
		return nil, errors.New("entry filename cannot be empty")
	}

	if len(oid) != 20 {
		return nil, fmt.Errorf("invalid SetOID length: expected 20 bytes, got %d", len(oid))
	}

	return &Entry{
		Name:        filename,
		AbsFilePath: absFilePath,
		OID:         oid,
		Executable:  executable,
	}, nil
}

func (e *Entry) GetRelativeFilePath(rootDir string) string {
	return strings.Replace(e.AbsFilePath, rootDir+"/", "", 1)
}

func (e *Entry) Content() ([]byte, error) {
	mode := regularMode
	if e.Executable {
		mode = executableMode
	}

	path := strings.Split(e.Name, "/")

	return []byte(fmt.Sprintf("%s %s\x00%s", mode, path[len(path)-1], e.OID)), nil
}
