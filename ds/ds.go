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
	validKey = regexp.MustCompile("[a-zA-Z](\\.?[0-9a-zA-Z_])*")
}

type DataSource interface {
	Get(key string) (string, error)
	Set(key, val string) error
	Value() (string, error)
	Watch() (chan struct{}, chan error)
	Feedback(error bool, message string) error
	Close() error
}

type Args struct {
	DS    string
	Token string
	Key   string
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
		return NewGRPC(m.Host, args.Token, args.Key)
	default:
		return nil, fmt.Errorf(`illegal datasource "%s"`, m.Scheme)
	}
}
