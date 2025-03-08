package command

import (
	"errors"
	"fmt"
	"io"

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
		return nil, fmt.Errorf("run commit cmd: %w", err)
	}

	ref, err := c.repository.Refs.CurrentRef()
	if err != nil {
		return nil, fmt.Errorf("read current ref: %w", err)
	}

	m := ""
	if commit.Parent == "" {
		m = "root-commit"
	}

	msg := fmt.Sprintf("[%s (%s) %s] %s", ref, m, commit.OID[0:7], commit.Message)

	return []byte(msg), nil
}

func (c *CommitCmd) Output(msg []byte, err error, stdout io.Writer) (int, error) {
	if err != nil {
		if errors.Is(err, repository.ErrNoFilesToCommit) {
			fmt.Fprintf(stdout, "%s", err.Error())

			return 1, nil
		}

		return 1, fmt.Errorf("commit cmd: %w", err)
	}

	fmt.Fprint(stdout, string(msg))

	return 0, nil
}
