package ds

import (
	"fmt"
	"net/url"
	"path"
)

type DataSource interface {
	Get(key string) (string, error)
	Set(key, val string) error
	JSON() ([]byte, error)
	Watch() (chan struct{}, chan error)
	Close() error
}

type Args struct {
	DS    string
	Token string
	Root  string
}

func Provider(args Args) (DataSource, error) {
	m, err := url.Parse(args.DS)
	if err != nil {
		return nil, err
	}

	switch m.Scheme {
	case "file":
		return NewFile(path.Join(m.Host, m.Path))
	default:
		return nil, fmt.Errorf(`illege datasource "%s"`, m.Scheme)
	}
}
