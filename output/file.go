package output

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"io/ioutil"
	"path"
)

type File struct {
	name string
}

func NewFile(name string) (*File, error) {
	if path.Ext(name) != ".json" {
		return nil, errors.New("only support json")
	}
	return &File{name: name}, nil
}

func (f *File) Put(data []byte) error {
	if !govalidator.IsJSON(string(data)) {
		return errors.New("illegal json")
	}

	return ioutil.WriteFile(f.name, data, 0660)
}
