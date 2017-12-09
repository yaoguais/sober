package authorize

import (
	"strings"
	"sync"
)

type Auth struct {
	Token string
	Path  string
}

var tokens sync.Map

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
