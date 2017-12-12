package dispatcher

import (
	"strings"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/store"
)

var (
	DefaultInterval = 5 * time.Second

	clients sync.Map
	locker  sync.Mutex
)

type client struct {
	id   string
	path string
	evtC chan store.Event
}

func (c *client) Event() <-chan store.Event {
	return c.evtC
}

func NewClient(path string) *client {
	return &client{
		id:   uuid.NewV4().String(),
		path: path,
		evtC: make(chan store.Event, 10),
	}
}

func Register(c *client) {
	locker.Lock()
	defer locker.Unlock()

	var cs []*client
	if v, ok := clients.Load(c.path); ok {
		cs = v.([]*client)
	}
	cs = append(cs, c)
	clients.Store(c.path, cs)
}

func UnRegister(c *client) {
	locker.Lock()
	defer locker.Unlock()

	if cs, ok := clients.Load(c.path); ok {
		ncs := []*client{}
		for _, v := range cs.([]*client) {
			if v.id != c.id {
				ncs = append(ncs, v)
			}
		}
		clients.Store(c.path, ncs)
	}
}

func Dispatch(s store.Store) {
	evtC, errC := s.Watch("")
	evtQueue := make(chan store.Event, 10000)

	go func() {
		for {
			m := make(map[string]store.Event)
			for len(evtQueue) > 0 {
				e := <-evtQueue
				m[e.Key] = e
			}

			for _, e := range m {
				clients.Range(func(key, val interface{}) bool {
					dispatchEvent(key, val, e)
					return true
				})
			}

			time.Sleep(DefaultInterval)
		}
	}()

	for {
		select {
		case err := <-errC:
			logrus.WithError(err).Error("dispatcher watch")
		case e := <-evtC:
			logrus.WithField("event", e).Debug("dispatch")
			evtQueue <- e
		}
	}
}

func dispatchEvent(key, val interface{}, e store.Event) {
	path := key.(string)
	if !strings.HasPrefix(e.Key, path) {
		return
	}

	cs, ok := val.([]*client)
	if !ok {
		return
	}

	for _, c := range cs {
		select {
		case c.evtC <- e:
			logrus.WithField("event", e).WithField("client", c).Debug("dispatcher")
		default:
		}
	}
}
