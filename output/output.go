package output

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"io/ioutil"
	"net/url"
	"path"
)

type Puter interface {
	Put([]byte) error
}

func Provider(output string) (Puter, error) {
	m, err := url.Parse(output)
	if err != nil {
		return nil, err
	}

	switch m.Scheme {
	case "file":
		return NewFileOutput(path.Join(m.Host, m.Path))
	default:
		return nil, fmt.Errorf(`illege output "%s"`, m.Scheme)
	}
}

type FileOutput struct {
	name string
}

func NewFileOutput(name string) (*FileOutput, error) {
	if path.Ext(name) != ".json" {
		return nil, errors.New("only support json")
	}
	return &FileOutput{name: name}, nil
}

func (f *FileOutput) Put(data []byte) error {
	if !govalidator.IsJSON(string(data)) {
		return errors.New("illegal json")
	}

	return ioutil.WriteFile(f.name, data, 0660)
}
