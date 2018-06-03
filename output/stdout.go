package output

import (
	"errors"
	"os"
)

type Stdout struct {
}

func NewStdout() (*Stdout, error) {
	return &Stdout{}, nil
}

func (*Stdout) Put(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty data")
	}

	os.Stdout.Write(data)

	return nil
}
