package output

import (
	"fmt"
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
		return NewFile(path.Join(m.Host, m.Path))
	case "stdout":
		return NewStdout()
	default:
		return nil, fmt.Errorf(`illegal output "%s"`, m.Scheme)
	}
}
