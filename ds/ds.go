package ds

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
)

var (
	validKey        *regexp.Regexp
	ErrIllegalKey   = errors.New("illegal key")
	ErrKeyNotExists = errors.New("key not exists")
)

func init() {
	validKey, _ = regexp.Compile("[a-zA-Z](\\.?[0-9a-zA-Z_])*")
}

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
	case "grpc":
		return NewGRPC(m.Host, args.Token, args.Root)
	default:
		return nil, fmt.Errorf(`illegal datasource "%s"`, m.Scheme)
	}
}
