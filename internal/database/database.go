package database

type Database struct {
	path string
}

type Object interface {
	Id() string
	Type() string
	String() string
}

// New
// objectsPath = .git/objects
func New(objectsPath string) *Database {
	return &Database{path: objectsPath}
}

func (d *Database) Store(t *Tree) error {

	t.
}


func (d *Database) writeObject(oid string, content []byte) error {

}
