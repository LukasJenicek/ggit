package database

type Commit struct{}

func NewCommit() *Commit {
	return &Commit{}
}

func (c Commit) Id() string {
	// TODO implement me
	panic("implement me")
}

func (c Commit) Type() string {
	// TODO implement me
	panic("implement me")
}

func (c Commit) String() string {
	// TODO implement me
	panic("implement me")
}
