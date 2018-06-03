package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	soberetry "github.com/yaoguais/sober/retry"
)

type Etcd struct {
	common
	kv   *clientv3.Client
	done chan struct{}
}

func NewEtcd(kv *clientv3.Client) (*Etcd, error) {
	return &Etcd{
		kv:   kv,
		done: make(chan struct{}),
	}, nil
}

func (s *Etcd) Get(key string) (string, error) {
	logrus.WithField("key", key).Debug("etcd get")

	resp, err := s.kv.Get(context.Background(), s.realKey(key))
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("key %s not found", key)
	}

	if len(resp.Kvs) != 1 {
		return "", errors.New("should get 1 kv pair")
	}

	return string(resp.Kvs[0].Value), nil
}

func (s *Etcd) Set(key, value string) error {
	_, err := s.kv.Put(context.Background(), s.realKey(key), value)

	logrus.WithField("key", key).Debug("etcd set")

	return err
}

func (s *Etcd) Watch(key string) (chan Event, chan error) {
	eventC := make(chan Event, 10)
	errC := make(chan error, 1)

	logrus.WithField("key", key).Debug("etcd watch")

	go func() {
		realKey := s.realKey(key)
		retry := soberetry.New(1, 60)

	try:
		c := s.kv.Watch(
			context.Background(),
			realKey,
			clientv3.WithPrefix())

		for {
			select {
			case <-s.done:
				return
			case resp := <-c:
				err := resp.Err()
				if err != nil {
					errC <- err
					retry.Wait()
					goto try
				}

				retry.Reset()

				for _, e := range resp.Events {
					logrus.WithField("event", e).Debug("etcd event")

					k := s.orignalKey(string(e.Kv.Key))
					if k == "" {
						continue
					}

					eventC <- Event{
						Key: k,
					}
				}

			}
		}
	}()

	return eventC, errC
}

func (s *Etcd) Close() error {
	if s.done == nil {
		return errors.New("close twice")
	}

	close(s.done)
	s.done = nil

	return nil
}
