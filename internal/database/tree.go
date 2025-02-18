package database

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	regularMode    = "100644"
	executableMode = "100755"
	directoryMode  = "40000"
)

type Entry struct {
	Name       string
	oid        []byte
	executable bool
}

func NewEntry(name string, oid []byte, executable bool) (*Entry, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("entry name cannot be empty")
	}

	if len(oid) != 20 {
		return nil, fmt.Errorf("invalid oid length: expected 20 bytes, got %d", len(oid))
	}

	return &Entry{
		Name:       name,
		oid:        oid,
		executable: executable,
	}, nil
}

func (e *Entry) Content() ([]byte, error) {
	mode := regularMode
	if e.executable {
		mode = executableMode
	}

	path := strings.Split(e.Name, "/")

	return []byte(fmt.Sprintf("%s %s\x00%s", mode, path[len(path)-1], e.oid)), nil
}

type Tree struct {
	// could be either Entry or Tree
	entries []Object
	name    string
	oid     []byte
	parent  *Tree
}

func NewTree(parent *Tree, name string) *Tree {
	return &Tree{
		parent: parent,
		name:   name,
	}
}

func (t *Tree) AddEntry(entry Object) {
	t.entries = append(t.entries, entry)
}

func (t *Tree) Content() ([]byte, error) {
	content := ""

	for _, entry := range t.entries {
		tree, ok := entry.(*Tree)
		if ok {
			if string(tree.oid) == "" {
				return nil, fmt.Errorf("tree %s has no oid", t.name)
			}

			content += fmt.Sprintf("%s %s\x00%s", directoryMode, tree.name, tree.oid)

			continue
		}

		contentPayload, err := entry.Content()
		if err != nil {
			return nil, fmt.Errorf("could not get content: %w", err)
		}

		content += string(contentPayload)
	}

	return []byte(fmt.Sprintf("tree %d\x00%s", len([]byte(content)), []byte(content))), nil
}

func findTrees(root *Tree) []*Tree {
	var trees []*Tree

	for _, e := range root.entries {
		t, ok := e.(*Tree)
		if ok {
			trees = append(trees, t)
			trees = append(trees, findTrees(t)...)
		}
	}

	return trees
}

//nolint:unused
func findLeaves(parent *Tree) []*Tree {
	var leaves []*Tree

	isLeafNode := true

	for i, e := range parent.entries {
		t, ok := e.(*Tree)
		if ok {
			isLeafNode = false

			leaves = append(leaves, findLeaves(t)...) // Merge results
		}

		if i == len(parent.entries)-1 && isLeafNode {
			leaves = append(leaves, parent) // Append the leaf node
		}
	}

	return leaves
}

// Build constructs a recursive tree structure from a sorted slice of entries.
// It maintains parent-child relationships allowing traversal from leaf nodes to the root.
// The entries must be sorted by name to ensure consistent tree building.
func Build(root *Tree, entries []*Entry) (*Tree, error) {
	treeCache := make(map[string]*Tree)

	for _, entry := range entries {
		d := filepath.Dir(entry.Name)
		// root folder
		if d == "." {
			root.AddEntry(entry)

			continue
		}

		t, ok := treeCache[d]
		if !ok {
			paths := strings.Split(d, "/")

			if len(paths) == 1 {
				t = NewTree(root, paths[0])
				treeCache[d] = t
				root.AddEntry(t)
			} else {
				parentPath := strings.Join(paths[:len(paths)-1], "/")

				parent, ok := treeCache[parentPath]
				if !ok {
					return nil, fmt.Errorf("parent entry %q not found", entry.Name)
				}

				t = NewTree(parent, paths[len(paths)-1])
				parent.AddEntry(t)
				treeCache[d] = t
			}
		}

		t.AddEntry(entry)
	}

	return root, nil
}
