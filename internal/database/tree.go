package database

import (
	"crypto/sha1"
	"fmt"
	"sort"
)

type blobs []*Blob

func (bs blobs) Len() int {
	return len(bs)
}

func (bs blobs) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

func (bs blobs) Less(i, j int) bool {
	return bs[i].filename < bs[j].filename
}

const (
	regularMode    = "100644"
	executableMode = "100755"
)

type Tree struct {
	blobs blobs
}

func NewTree(blobs []*Blob) *Tree {
	return &Tree{blobs: blobs}
}

func (t *Tree) ID() []byte {
	hasher := sha1.New()
	hasher.Write(t.Content())

	return hasher.Sum(nil)
}

func (t *Tree) Type() string {
	return "tree"
}

func (t *Tree) Content() []byte {
	sort.Sort(t.blobs)

	content := ""

	for _, blob := range t.blobs {
		mode := regularMode
		if blob.isExecutable {
			mode = executableMode
		}

		content += fmt.Sprintf("%s %s\x00%s", mode, blob.filename, blob.ID())
	}

	return []byte(fmt.Sprintf("tree %d\x00%s", len([]byte(content)), []byte(content)))
}
