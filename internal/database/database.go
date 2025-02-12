package database

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

type Database struct {
	gitRootDir  string
	objectsPath string
}

type Object interface {
	ID() []byte
	Type() string
	Content() []byte
}

func New(gitDir string) *Database {
	return &Database{
		gitRootDir:  gitDir,
		objectsPath: gitDir + "/objects",
	}
}

func (d *Database) Store(o Object) error {
	c := o.Content()
	oid := o.ID()

	if len(oid) != 20 {
		return fmt.Errorf("sha1 hash must have 20 bytes, invalid object id: %s", oid)
	}

	return d.writeObject(hex.EncodeToString(oid), c)
}

func (d *Database) writeObject(oid string, content []byte) error {
	realPath := fmt.Sprintf("%s/%s/%s", d.objectsPath, oid[:2], oid[2:])
	dir := filepath.Dir(realPath)

	// create directory in .git/objects folder
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("could not create directory %s: %w", dir, err)
			}
		} else {
			return fmt.Errorf("check directory: %w", err)
		}
	}

	// object already exist do not overwrite
	if _, err := os.Stat(realPath); err == nil {
		return nil
	}

	// zlib compression with best speed
	var compressed bytes.Buffer

	zlibWriter, err := zlib.NewWriterLevel(&compressed, zlib.BestSpeed)
	if err != nil {
		return fmt.Errorf("could not create zlib writer: %w", err)
	}

	if _, err := zlibWriter.Write(content); err != nil {
		return fmt.Errorf("compress object content: %w", err)
	}

	if err := zlibWriter.Close(); err != nil {
		return fmt.Errorf("close zlib writer: %w", err)
	}

	// when os write content to file it does not have to be done at once
	// first we create tmp object and then move it ( rename it )
	tmpPath := fmt.Sprintf("%s/%s/tmp_%s", d.objectsPath, oid[:2], oid[2:])
	if err := os.WriteFile(tmpPath, compressed.Bytes(), 0o644); err != nil {
		return fmt.Errorf("store tmp object %s: %w", oid, err)
	}

	if err = os.Rename(tmpPath, realPath); err != nil {
		return fmt.Errorf("rename tmp object %s: %w", oid, err)
	}

	return nil
}
