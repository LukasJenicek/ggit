package database

import "github.com/LukasJenicek/ggit/internal/workspace"

const mode = "100644"

type Tree struct {
	files []*workspace.File
}

func NewTree(files []*workspace.File) *Tree {
	return &Tree{files: files}
}

func (tree *Tree) Id() string {
	return ""
}

func (tree *Tree) Type() string {
	return "tree"
}

func (tree *Tree) String() string {

}
