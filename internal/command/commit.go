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
	cID, err := c.repository.Commit()
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("[%s] Successfully committed changes\n", cID)), nil
}
