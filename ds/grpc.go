package ds

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/client"
	"github.com/yaoguais/sober/kvpb"
	"google.golang.org/grpc"
)

type GRPC struct {
	kv   *client.KV
	key  string
	data string
	done chan struct{}
	sync.RWMutex
}

func NewGRPC(addr, token, key string) (*GRPC, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	kvc := kvpb.NewKVClient(conn)

	kv := client.NewKV(token, key, kvc)

	g := &GRPC{
		kv:   kv,
		key:  key,
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
	return "", errors.New("not support")
}

func (g *GRPC) Set(key, val string) error {
	return errors.New("forbid set")
}

func (g *GRPC) Value() (string, error) {
	g.RLock()
	v := g.data
	g.RUnlock()

	return v, nil
}

func (g *GRPC) Watch() (chan struct{}, chan error) {
	c := make(chan struct{}, 1)

	evtC, errC := g.kv.Watch()

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

func (g *GRPC) Feedback(error bool, message string) error {
	return g.kv.Feedback(error, message)
}

func (g *GRPC) Close() error {
	if g.done == nil {
		return errors.New("close twice")
	}

	close(g.done)
	g.done = nil

	return nil
}

func (g *GRPC) load() (string, error) {
	v, err := g.kv.Get()
	if err != nil {
		return "", err
	}

	logrus.WithField("value", v).Debug("ds load")

	return v, nil
}
