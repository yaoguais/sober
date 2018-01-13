package authorize

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/dispatcher"
	"github.com/yaoguais/sober/store"
	"strings"
	"sync"
)

type Auth struct {
	Name  string
	Token string
	Path  string
}

var (
	tokens *sync.Map
)

func replaceAll(auths []Auth) {
	logrus.Debug("replace all authorize")
	tmpTokens := &sync.Map{}
	for _, v := range auths {
		tmpTokens.Store(v.Token, v)
		logrus.WithField("auth", v).Debug("add auth")
	}
	tokens = tmpTokens
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

	if auths, err := load(s, path); err != nil {
		return err
	} else {
		replaceAll(auths)
	}

	c := dispatcher.NewClient(path)
	dispatcher.Register(c)

	go func() {
		defer dispatcher.UnRegister(c)
		for range c.Event() {
			if auths, err := load(s, path); err != nil {
				logrus.WithError(err).Error("load auth")
			} else {
				replaceAll(auths)
			}
		}
	}()

	return nil
}

func load(s store.Store, path string) ([]Auth, error) {
	kv, err := s.KV(path)
	if err != nil {
		return nil, err
	}

	var auths []Auth
	for k, v := range kv {
		if v == "" {
			continue
		}
		keys := strings.Split(k, "/")
		if len(keys) == 3 {
			name := keys[1]
			if keys[2] == "path" {
				tokenKey := fmt.Sprintf("/%s/token", name)
				if token := kv[tokenKey]; token != "" {
					auth := Auth{
						Name:  name,
						Token: token,
						Path:  v,
					}
					auths = append(auths, auth)
				}
			}
		}
	}

	if err := check(auths); err != nil {
		return nil, err
	}

	return auths, nil
}

func check(auths []Auth) error {
	m := make(map[string]struct{})
	for _, v := range auths {
		if _, ok := m[v.Token]; ok {
			return fmt.Errorf("duplicate token %s", v.Token)
		}
		m[v.Token] = struct{}{}
	}
	return nil
}
