package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/yaoguais/sober/authorize"
	"regexp"
	"strings"
	"time"
)

type Etcd struct {
	kv      *clientv3.Client
	root    string
	rule    *regexp.Regexp
	svcPath string
	c       chan struct{}
}

func NewEtcd(root, rule, svcPath string, options clientv3.Config) (*Etcd, error) {
	reg, err := regexp.Compile(rule)
	if err != nil {
		return nil, err
	}

	kv, err := clientv3.New(options)
	if err != nil {
		return nil, err
	}

	s := &Etcd{
		root:    root,
		rule:    reg,
		svcPath: svcPath,
		kv:      kv,
		c:       make(chan struct{}),
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	if err := s.watch(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Etcd) load() error {
	resp, err := s.kv.Get(
		context.Background(),
		s.authorizePath(),
		clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, v := range resp.Kvs {
		auth := authorize.Auth{
			Token: string(v.Key),
			Path:  string(v.Value),
		}
		authorize.Add(auth)
	}

	return nil
}

func (s *Etcd) watch() error {
	go func() {
		init := 1 * time.Second
		max := time.Minute
		sleep := init

	retry:
		c := s.kv.Watch(
			context.Background(),
			s.authorizePath(),
			clientv3.WithPrefix())

		for {
			select {
			case <-s.c:
				return
			case resp := <-c:
				if resp.Err() != nil {
					sleep *= 2
					if sleep > max {
						sleep = max
					}
					time.Sleep(sleep)
					goto retry
				}

				sleep = init

				for _, e := range resp.Events {
					fmt.Printf("e: %v\n", *e)
					auth := authorize.Auth{
						Token: string(e.Kv.Key),
						Path:  string(e.Kv.Value),
					}
					switch e.Type {
					case clientv3.EventTypePut:
						authorize.Add(auth)
					case clientv3.EventTypeDelete:
						authorize.Remove(auth)
					}
				}
			}
		}
	}()

	return nil
}

func (s *Etcd) path(path string) string {
	return strings.ToLower(s.root + path)
}

func (s *Etcd) authorizePath() string {
	return s.path(s.svcPath + "/authorize")
}

func (s *Etcd) Close() error {
	if s.c == nil {
		return errors.New("close twice")
	}

	close(s.c)
	s.c = nil

	return nil
}
