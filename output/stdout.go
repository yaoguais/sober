package output

import (
	"errors"
	"github.com/asaskevich/govalidator"
)

type Stdout struct {
}

func NewStdout() (*Stdout, error) {
	return &Stdout{}, nil
}

func (*Stdout) Put(data []byte) error {
	if !govalidator.IsJSON(string(data)) {
		return errors.New("illegal json")
	}

	return nil
}
