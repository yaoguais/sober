package authorize

import (
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/store"
	"strings"
	"sync"
)

type Auth struct {
	Token string
	Path  string
}

var (
	tokens sync.Map
)

func Add(auth Auth) {
	tokens.Store(auth.Token, auth)
	logrus.WithField("auth", auth).Debug("add auth")
}

func Remove(auth Auth) {
	tokens.Delete(auth.Token)
	logrus.WithField("auth", auth).Debug("remove auth")
}

func Range(fn func(Auth)) {
	tokens.Range(func(key, val interface{}) bool {
		fn(val.(Auth))
		return true
	})
}

func Valid(check Auth) bool {
	if v, ok := tokens.Load(check.Token); !ok {
		return false
	} else {
		vv := v.(Auth)
		return strings.HasPrefix(check.Path, vv.Path)
	}
}

func Start(s store.Store, basePath string) error {
	path := basePath + "/authorize"

	kv, err := s.KV(path)
	if err != nil {
		return err
	}

	for k, v := range kv {
		auth := Auth{
			Token: k[1:],
			Path:  v,
		}
		Add(auth)
	}

	evtC, errC := s.Watch(path)

	go func() {
		select {
		case <-errC:
		case evts := <-evtC:
			for _, e := range evts {
				auth := Auth{
					Token: string(e.Key)[1:],
					Path:  string(e.Value),
				}
				switch e.Type {
				case store.EventTypePut:
					Add(auth)
				case store.EventTypeDelete:
					Remove(auth)
				}
			}
		}
	}()

	return nil
}
