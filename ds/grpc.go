package ds

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/client"
	soberini "github.com/yaoguais/sober/ini"
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
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	kvc := kvpb.NewKVClient(conn)
	kv.SetKVC(kvc)

	g := &GRPC{
		kv:   kv,
		root: root,
		done: make(chan struct{}),
	}

	data, err := g.load()
	if err != nil {
		return nil, err
	}

	g.data = data

	return g, nil
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

	return soberini.IniToPrettyJSON(g.data)
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
				data, err := g.load()
				if err != nil {
					errC <- err
				} else {
					g.Lock()
					g.data = data
					g.Unlock()
					c <- struct{}{}
					logrus.WithField("event", e).Debug("ds event")
				}
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

func (g *GRPC) load() (map[string]string, error) {
	m, err := g.kv.Get(g.root)
	if err != nil {
		return nil, err
	}

	logrus.WithField("kv", m).Debug("ds load")

	data := make(map[string]string)
	for k, v := range m {
		k := strings.TrimLeft(k, "/")
		k = strings.Replace(k, "/", ".", -1)
		data[k] = v
	}

	return data, nil
}
