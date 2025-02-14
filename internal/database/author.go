package database

import (
	"fmt"
	"time"
)

type Author struct {
	Email string
	Name  string
	Now   *time.Time
}

func NewAuthor(email, name string, now *time.Time) *Author {
	return &Author{
		Email: email,
		Name:  name,
		Now:   now,
	}
}

// Name <email> unix-timestamp timezone.
func (a *Author) String() string {
	return fmt.Sprintf(
		"%s <%s> %d %s",
		a.Name,
		a.Email,
		a.Now.Unix(),
		a.Now.Format("-0700"),
	)
}
