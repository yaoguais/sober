package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/yaoguais/sober"
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
		kv[k] = string(v.Value)
	}

	return kv, nil
}

func (s *Etcd) Watch(path string) (chan []Event, chan error) {
	errC := make(chan error, 1)
	realPath := s.realPath(path)
	retry := sober.NewRetry(1, 60)
	eventC := make(chan []Event, 10)

	fmt.Printf("realPath:%v\n", realPath)

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

				fmt.Printf("events:%v\n", resp.Events)

				events := make([]Event, len(resp.Events))
				for _, e := range resp.Events {
					k := s.orignalPath(string(e.Kv.Key))
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
						Value: string(e.Kv.Value),
					}
					events = append(events, evt)
				}

				eventC <- events
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
