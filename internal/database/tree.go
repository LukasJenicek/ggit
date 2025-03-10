package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	regularMode    = "100644"
	executableMode = "100755"
	directoryMode  = "40000"
)

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

func NewRootTree() *Tree {
	return &Tree{
		parent: nil,
		name:   "",
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
				return nil, fmt.Errorf("tree %s has no SetOID", t.name)
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
		d := filepath.Dir(entry.AbsFilePath)
		// root folder
		if d == "." {
			root.AddEntry(entry)

			continue
		}

		t, ok := treeCache[d]
		if !ok {
			split := strings.Split(d, string(os.PathSeparator))

			for i := range split {
				folderParentPath := strings.Join(split[:i], string(os.PathSeparator))
				folderPath := strings.Join(split[:i+1], string(os.PathSeparator))

				if folderParentPath == "" {
					treeCache[split[0]] = NewTree(root, split[0])
					root.AddEntry(treeCache[split[0]])

					continue
				}

				parent, ok := treeCache[folderParentPath]
				if !ok {
					return nil, fmt.Errorf("tree %s does not exist", folderParentPath)
				}

				t, ok = treeCache[folderPath]
				if !ok {
					t = NewTree(parent, split[i])
					treeCache[split[i]] = t
				}

				parent.AddEntry(t)
			}
		}

		t.AddEntry(entry)
	}

	return root, nil
}
