package database

import "os"

type Blob struct {
}

func NewBlob(file *os.File) *Blob {

}

func (b Blob) Id() string {
	//TODO implement me
	panic("implement me")
}

func (b Blob) Type() string {
	return "blob"
}

func (b Blob) String() string {
}
