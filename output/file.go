package output

import (
	"errors"
	"io/ioutil"
)

type File struct {
	name string
}

func NewFile(name string) (*File, error) {
	return &File{name: name}, nil
}

func (f *File) Put(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty data")
	}

	return ioutil.WriteFile(f.name, data, 0660)
}
