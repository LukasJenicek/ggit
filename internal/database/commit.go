package database

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"
)

type Commit struct {
	treeOid string
	author  *Author
	message string
	parent  string
}

func NewCommit(treeOid string, author *Author, message string, parent string) (*Commit, error) {
	if strings.TrimSpace(treeOid) == "" {
		return nil, errors.New("treeOid must not be empty")
	}

	if strings.TrimSpace(message) == "" {
		return nil, errors.New("commit message cannot be empty")
	}

	return &Commit{
		treeOid: treeOid,
		author:  author,
		message: message,
		parent:  parent,
	}, nil
}

func (c *Commit) ID() []byte {
	hasher := sha1.New()
	hasher.Write(c.Content())

	return hasher.Sum(nil)
}

func (c *Commit) Type() string {
	return "commit"
}

func (c *Commit) Content() []byte {
	lines := []string{"tree " + c.treeOid}
	if c.parent != "" {
		lines = append(lines, "parent "+c.parent)
	}
	lines = append(lines, "author "+c.author.String(), "committer "+c.author.String())
	lines = append(lines, "", c.message)

	content := strings.Join(lines, "\n")
	content += "\n"

	return []byte(fmt.Sprintf("%s %d\x00%s", c.Type(), len(content), content))
}
