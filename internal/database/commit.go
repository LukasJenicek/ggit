package database

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

type Commit struct {
	treeOid string
	author  *Author
	message string
}

func NewCommit(treeOid string, author *Author, message string) *Commit {
	return &Commit{
		treeOid: treeOid,
		author:  author,
		message: message,
	}
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
	lines := []string{
		"tree " + c.treeOid,
		"author " + c.author.String(),
		"committer " + c.author.String(),
		"",
		c.message,
	}

	content := strings.Join(lines, "\n")
	content += "\n"

	return []byte(fmt.Sprintf("%s %d\x00%s", c.Type(), len(content), content))
}
