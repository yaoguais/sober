package dispatcher

import (
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/store"
	"strings"
	"sync"
	"time"
)

var (
	clients         sync.Map
	DefaultInterval = 5 * time.Second
	locker          sync.Mutex
)

type client struct {
	id   string
	path string
	push time.Time
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

	for {
		select {
		case err := <-errC:
			logrus.WithError(err).Error("dispatcher watch")
		case e := <-evtC:
			logrus.WithField("event", e).Debug("dispatcher")

			clients.Range(func(key, val interface{}) bool {
				dispatchEvent(key, val, e)
				return true
			})
		}
	}
}

func dispatchEvent(key, val interface{}, e store.Event) {
	now := time.Now()
	path := key.(string)
	if strings.HasPrefix(e.Key, path) {
		if cs, ok := val.([]*client); ok {
			for _, c := range cs {
				if now.Sub(c.push) > DefaultInterval {
					c.push = now
					select {
					case c.evtC <- e:
						logrus.WithField("event", e).
							WithField("client", c).
							Debug("dispatcher")
					default:
					}
				}
			}
		}
	}
}
