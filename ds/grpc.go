package ds

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober"
	"github.com/yaoguais/sober/client"
	"github.com/yaoguais/sober/kvpb"
	"google.golang.org/grpc"
	"strings"
	"sync"
)

type GRPC struct {
	kv   *client.KV
	root string
	data map[string]string
	done chan struct{}
	sync.RWMutex
}

func NewGRPC(addr, token, root string) (*GRPC, error) {
	kv := client.NewKV(token, root)
	conn, err := grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(kv))
	if err != nil {
		return nil, err
	}

	kvc := kvpb.NewKVClient(conn)
	kv.SetKVC(kvc)

	m, err := kv.Get(root)
	if err != nil {
		return nil, err
	}

	logrus.WithField("kv", m).Debug("ds fetch")

	data := make(map[string]string)
	for k, v := range m {
		k := strings.TrimLeft(k, "/")
		k = strings.Replace(k, "/", ".", -1)
		data[k] = v
	}

	return &GRPC{
		kv:   kv,
		root: root,
		data: data,
		done: make(chan struct{}),
	}, nil
}

func (g *GRPC) Get(key string) (string, error) {
	g.RLock()
	defer g.RUnlock()

	if !validKey.Match([]byte(key)) {
		return "", ErrIllegalKey
	}

	if v, ok := g.data[key]; !ok {
		return "", ErrKeyNotExists
	} else {
		return v, nil
	}
}

func (g *GRPC) Set(key, val string) error {
	return errors.New("forbid set")
}

func (g *GRPC) JSON() ([]byte, error) {
	g.RLock()
	defer g.RUnlock()

	if len(g.data) == 0 {
		return nil, errors.New("empty data")
	}

	return sober.IniToPrettyJSON(g.data)
}

func (g *GRPC) Watch() (chan struct{}, chan error) {
	c := make(chan struct{}, 1)

	evtC, errC := g.kv.Watch(g.root)

	go func() {
		for {
			select {
			case <-g.done:
				g.kv.Cancel()
				return
			case e := <-evtC:
				c <- struct{}{}
				logrus.WithField("event", e).Debug("ds event")
			}
		}
	}()

	return c, errC
}

func (g *GRPC) Close() error {
	if g.done == nil {
		return errors.New("close twice")
	}

	close(g.done)
	g.done = nil

	return nil
}
