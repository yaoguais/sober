package store

import (
	"context"
	"errors"
	gopath "path"

	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/crypto"
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

func (s *Etcd) KV(path string) (map[string]string, error) {
	realPath := s.realPath(path)
	resp, err := s.kv.Get(
		context.Background(),
		realPath,
		clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	kv := make(map[string]string)
	rl := len(realPath)
	for _, v := range resp.Kvs {
		k := string(v.Key)[rl:]
		if k == "" {
			continue
		}
		kv[k] = crypto.Decode(gopath.Base(k), string(v.Value))
	}

	return kv, nil
}

func (s *Etcd) Set(path string, kv map[string]string) error {
	pkv, err := s.KV(path)
	if err != nil {
		return err
	}

	realPath := s.realPath(path)
	var ops []clientv3.Op
	for k, v := range kv {
		key := gopath.Join(realPath, k)
		op := clientv3.OpPut(key, crypto.Encode(gopath.Base(key), v))
		ops = append(ops, op)
	}
	for k := range pkv {
		if _, ok := kv[k]; !ok {
			key := gopath.Join(realPath, k)
			ops = append(ops, clientv3.OpDelete(key))
		}
	}

	ctx := context.Background()
	txn := s.kv.Txn(ctx)
	_, err = txn.Then(ops...).Commit()

	return err
}

func (s *Etcd) Watch(path string) (chan Event, chan error) {
	eventC := make(chan Event, 10)
	errC := make(chan error, 1)

	realPath := s.realPath(path)
	pathLen := len(path)

	retry := soberetry.New(1, 60)

	logrus.WithField("path", realPath).Debug("etcd watch")

	go func() {
	try:
		c := s.kv.Watch(
			context.Background(),
			realPath,
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

					k := s.orignalPath(string(e.Kv.Key))[pathLen:]
					if k == "" {
						continue
					}

					var t EventType
					switch e.Type {
					case clientv3.EventTypePut:
						t = EventTypePut
					case clientv3.EventTypeDelete:
						t = EventTypeDelete
					}
					evt := Event{
						Type:  t,
						Key:   k,
						Value: crypto.Decode(gopath.Base(k), string(e.Kv.Value)),
					}

					eventC <- evt
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
