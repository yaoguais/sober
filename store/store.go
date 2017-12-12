package store

import (
	"github.com/yaoguais/sober/kvpb"
)

type EventType = kvpb.Event_EventType

const (
	EventTypePut    = kvpb.Event_PUT
	EventTypeDelete = kvpb.Event_DELETE
)

type Event struct {
	Type  EventType
	Key   string
	Value string
}

type Store interface {
	KV(path string) (map[string]string, error)
	Watch(path string) (chan Event, chan error)
	Close() error
}

type common struct {
	root string
}

func (c *common) SetRoot(root string) {
	c.root = root
}

func (c *common) realPath(path string) string {
	return c.root + path
}

func (c *common) orignalPath(path string) string {
	return path[len(c.root):]
}

func (c *common) KV(path string) (map[string]string, error) {
	panic("abstract method")
	return nil, nil
}

func (c *common) Watch(path string) (chan Event, chan error) {
	panic("abstract method")
	return nil, nil
}

func (c *common) Close() error {
	panic("abstract method")
	return nil
}
