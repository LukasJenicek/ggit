package database

import (
	"fmt"
	"path/filepath"
	"strings"
)

const regularMode = "100644"
const executableMode = "100755"
const directoryMode = "40000"

type Tree struct {
	// could be either Entry or Tree
	entries []Object
}

// Build
// Recursive tree out of sorted entries slice
func Build(root *Tree, entries []*Entry) *Tree {
	treeMapCache := make(map[string]*Tree)

	for _, entry := range entries {
		d := filepath.Dir(entry.Name)
		// root folder
		if d == "." {
			root.AddEntry(entry)
			continue
		}

		t, ok := treeMapCache[d]
		if !ok {
			t = NewTree()
			treeMapCache[d] = t

			paths := strings.Split(d, "/")
			if len(paths) == 1 {
				root.AddEntry(t)
			} else {
				tPath := strings.Join(paths[:len(paths)-1], "/")
				r, ok := treeMapCache[tPath]
				if !ok {
					panic(fmt.Sprintf("tree map entry %q not found", entry.Name))
				}
				r.AddEntry(t)
			}
		}
		t.AddEntry(entry)
	}

	return root
}

type Entry struct {
	Name       string
	oid        []byte
	executable bool
}

func NewEntry(name string, oid []byte, executable bool) *Entry {
	return &Entry{
		Name:       name,
		oid:        oid,
		executable: executable,
	}
}

func (e *Entry) Content() []byte {
	mode := regularMode
	if e.executable {
		mode = executableMode
	}

	return []byte(fmt.Sprintf("%s %s\x00%s", mode, e.Name, e.oid))
}

func NewTree() *Tree { return &Tree{} }

func (t *Tree) AddEntry(entry Object) {
	t.entries = append(t.entries, entry)
}

func (t *Tree) Content() []byte {
	content := ""
	for _, entry := range t.entries {
		content += string(entry.Content())
	}

	return []byte(fmt.Sprintf("tree %d\x00%s", len([]byte(content)), []byte(content)))
}
