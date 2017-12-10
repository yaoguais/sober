package authorize

import (
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
}

func Remove(auth Auth) {
	tokens.Delete(auth.Token)
}

func Replace(auths []Auth) {
	t := sync.Map{}
	for _, v := range auths {
		t.Store(v.Token, v)
	}
	tokens = t
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
		return strings.HasPrefix(strings.ToLower(vv.Path), strings.ToLower(check.Path))
	}
}

func Start(s store.Store, basePath string) error {
	path := basePath + "/authorize"
	preLen := len(basePath) + 1

	kv, err := s.KV(path)
	if err != nil {
		return err
	}

	for k, v := range kv {
		auth := Auth{
			Token: k[preLen:],
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
					Token: string(e.Key),
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
