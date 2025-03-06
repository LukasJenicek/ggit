package database

import (
	"errors"
	"fmt"
	"strings"
)

type Commit struct {
	RootOID string
	OID     string
	Author  *Author
	Message string
	Parent  string
}

func NewCommit(parent string, rootOID string, author *Author, message string) (*Commit, error) {
	if strings.TrimSpace(rootOID) == "" {
		return nil, errors.New("root oid must not be empty")
	}

	if strings.TrimSpace(message) == "" {
		return nil, errors.New("commit message cannot be empty")
	}

	return &Commit{
		RootOID: rootOID,
		Author:  author,
		Message: message,
		Parent:  parent,
	}, nil
}

func (c *Commit) SetOID(oid string) error {
	if len(oid) != 40 {
		return fmt.Errorf("oid must be 40 characters long: %s", oid)
	}

	c.OID = oid

	return nil
}

func (c *Commit) Content() ([]byte, error) {
	lines := []string{"tree " + c.RootOID}
	if c.Parent != "" {
		lines = append(lines, "parent "+c.Parent)
	}

	lines = append(lines, "author "+c.Author.String(), "committer "+c.Author.String())
	lines = append(lines, "", c.Message)

	content := strings.Join(lines, "\n")
	content += "\n"

	return []byte(fmt.Sprintf("%s %d\x00%s", "commit", len(content), content)), nil
}
