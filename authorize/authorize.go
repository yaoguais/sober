package authorize

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/decode"
	"github.com/yaoguais/sober/dispatcher"
	"github.com/yaoguais/sober/store"
)

type Auth struct {
	Token string `json:"token" toml:"token" yaml:"toml"`
	Key   string `json:"key" toml:"key" yaml:"key"`
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
		return strings.HasPrefix(check.Key, vv.Key)
	}
}

func Start(s store.Store, key string) error {
	if auths, err := load(s, key); err != nil {
		return err
	} else {
		replaceAll(auths)
	}

	c := dispatcher.NewClient(key)
	dispatcher.Register(c)

	go func() {
		defer dispatcher.UnRegister(c)
		for range c.Event() {
			if auths, err := load(s, key); err != nil {
				logrus.WithError(err).Error("load auth")
			} else {
				replaceAll(auths)
			}
		}
	}()

	return nil
}

func load(s store.Store, key string) ([]Auth, error) {
	v, err := s.Get(key)
	if err != nil {
		return nil, err
	}

	var auths []Auth
	if err = decode.Decode(key, v, &auths); err != nil {
		return nil, err
	}

	if err = check(auths); err != nil {
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
