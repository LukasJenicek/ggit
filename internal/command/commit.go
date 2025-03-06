package command

import (
	"fmt"

	"github.com/LukasJenicek/ggit/internal/repository"
)

type CommitCmd struct {
	repository *repository.Repository
}

func NewCommitCmd(repository *repository.Repository) (*CommitCmd, error) {
	return &CommitCmd{repository: repository}, nil
}

func (c *CommitCmd) Run() ([]byte, error) {
	commit, err := c.repository.Commit()
	if err != nil {
		return nil, err
	}

	ref, err := c.repository.Refs.CurrentRef()
	if err != nil {
		return nil, fmt.Errorf("read current ref: %w", err)
	}

	msg := fmt.Sprintf("[%s] (root-commit) %s] %s", ref, commit.OID[0:7], commit.Message)

	return []byte(fmt.Sprintf(msg)), nil
}
