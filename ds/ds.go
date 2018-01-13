package ds

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
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

func ReplaceToDotKey(m map[string]string) map[string]string {
	data := make(map[string]string)
	for k, v := range m {
		k := strings.TrimLeft(k, "/")
		k = strings.Replace(k, "/", ".", -1)
		data[k] = v
	}
	return data
}

func ReplaceToSlashKey(m map[string]string) map[string]string {
	data := make(map[string]string)
	for k, v := range m {
		k = strings.Replace(k, ".", "/", -1)
		if len(k) > 0 && k[0] != '/' {
			k = "/" + k
		}
		data[k] = v
	}
	return data
}
