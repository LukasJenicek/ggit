package database

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LukasJenicek/ggit/internal/filesystem"
	"github.com/LukasJenicek/ggit/internal/hasher"
	"github.com/LukasJenicek/ggit/internal/workspace"
)

type Database struct {
	fs          filesystem.Fs
	gitRootDir  string
	objectsPath string
}

type Object interface {
	Content() ([]byte, error)
}

func New(fs filesystem.Fs, gitDir string) *Database {
	return &Database{
		fs:          fs,
		gitRootDir:  gitDir,
		objectsPath: gitDir + "/objects",
	}
}

// StoreTree
// I just process the slice in reverse order.
func (d *Database) StoreTree(root *Tree) ([]byte, error) {
	if root == nil {
		return nil, errors.New("root tree must be provided")
	}

	// TODO: Use queue and reverse the order ? Would be more memory efficient
	trees := findTrees(root)
	// start from end
	for i := len(trees) - 1; i >= 0; i-- {
		// TODO: detect circular reference ???
		t := trees[i]

		oid, err := d.Store(t)
		if err != nil {
			return nil, fmt.Errorf("store tree: %w", err)
		}

		t.oid = oid
	}

	oid, err := d.Store(root)
	if err != nil {
		return nil, fmt.Errorf("store root: %w", err)
	}

	return oid, nil
}

func (d *Database) Store(o Object) ([]byte, error) {
	c, err := o.Content()
	if err != nil {
		return nil, fmt.Errorf("get object content: %w", err)
	}

	oid, err := hasher.SHA1HashContent(c)
	if err != nil {
		return nil, fmt.Errorf("generate oid: %w", err)
	}

	if len(oid) != 20 {
		return nil, fmt.Errorf("invalid object id: %s, sha1 hash must have 20 bytes", oid)
	}

	err = d.writeObject(hex.EncodeToString(oid), c)
	if err != nil {
		return nil, fmt.Errorf("store object: %w", err)
	}

	return oid, nil
}

func (d *Database) SaveBlobs(files workspace.Set) ([]*Entry, error) {
	entries := make([]*Entry, 0, len(files))

	for _, file := range files.SortedValues() {
		entry, err := d.saveBlob(file)
		if err != nil {
			return nil, fmt.Errorf("save blob: %w", err)
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (d *Database) saveBlob(filePath string) (*Entry, error) {
	content, err := d.fs.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file content: %w", err)
	}

	oid, err := d.Store(NewBlob(content))
	if err != nil {
		return nil, fmt.Errorf("store blob %s: %w", filePath, err)
	}

	f, err := d.fs.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("stat file %s: %w", filePath, err)
	}

	executable := false
	if f.Mode().Perm()&0o100 != 0 {
		executable = true
	}

	return NewEntry(f.Name(), filePath, oid, executable)
}

func (d *Database) writeObject(oid string, content []byte) error {
	realPath := fmt.Sprintf("%s/%s/%s", d.objectsPath, oid[:2], oid[2:])
	dir := filepath.Dir(realPath)

	// object already exist do not overwrite
	if _, err := d.fs.Stat(realPath); err == nil {
		return nil
	}

	// create directory in .git/objects folder
	if _, err := d.fs.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := d.fs.Mkdir(dir, 0o755); err != nil {
				return fmt.Errorf("create objects directory %s: %w", dir, err)
			}
		} else {
			return fmt.Errorf("check directory: %w", err)
		}
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
	if err := d.fs.WriteFile(tmpPath, compressed.Bytes(), 0o644); err != nil {
		return fmt.Errorf("store tmp object %s: %w", tmpPath, err)
	}

	if err = d.fs.Rename(tmpPath, realPath); err != nil {
		return fmt.Errorf("rename tmp object %s: %w", oid, err)
	}

	return nil
}
