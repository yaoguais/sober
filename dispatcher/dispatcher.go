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
	key  string
	evtC chan store.Event
}

func (c *client) Event() <-chan store.Event {
	return c.evtC
}

func NewClient(key string) *client {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return &client{
		id:   id.String(),
		key:  key,
		evtC: make(chan store.Event, 10),
	}
}

func Register(c *client) {
	locker.Lock()
	defer locker.Unlock()

	var cs []*client
	if v, ok := clients.Load(c.key); ok {
		cs = v.([]*client)
	}
	cs = append(cs, c)
	clients.Store(c.key, cs)
}

func UnRegister(c *client) {
	locker.Lock()
	defer locker.Unlock()

	if cs, ok := clients.Load(c.key); ok {
		ncs := []*client{}
		for _, v := range cs.([]*client) {
			if v.id != c.id {
				ncs = append(ncs, v)
			}
		}
		clients.Store(c.key, ncs)
	}
}

func Dispatch(s store.Store) {
	evtC, errC := s.Watch("")
	evtQueue := make(chan store.Event, 10000)

	go func() {
		for {
			cc := make(map[string]struct{})
			for len(evtQueue) > 0 {
				e := <-evtQueue
				clients.Range(func(key, val interface{}) bool {
					for _, c := range val.([]*client) {
						if _, ok := cc[c.id]; !ok {
							cc[c.id] = struct{}{}

							k := key.(string)
							if strings.HasPrefix(e.Key, k) {
								select {
								case c.evtC <- e:
									logrus.WithField("event", e).WithField("client", c).Debug("dispatcher")
								default:
								}
							}
						}
					}
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
