package output

import (
	"errors"
	"os"

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

	os.Stdout.Write(data)

	return nil
}
