package store

import (
	"errors"
	"regexp"
	"strings"
)

type EventType int

const (
	EventTypePut    = 0
	EventTypeDelete = 1
)

var (
	ErrIllegalPath = errors.New("illegal path")
)

type Event struct {
	Type  EventType
	Key   string
	Value string
}

type Store interface {
	KV(path string) (map[string]string, error)
	Watch(path string) (chan []Event, chan error)
	Close() error
}

type common struct {
	root string
	rule *regexp.Regexp
}

func (c *common) SetRoot(root string) *common {
	c.root = strings.ToLower(root)
	return c
}

func (c *common) SetRule(rule *regexp.Regexp) *common {
	c.rule = rule
	return c
}

func (c *common) ValidPath(path string) bool {
	return c.rule.Match([]byte(path))
}

func (c *common) realPath(path string) string {
	return c.root + strings.ToLower(path)
}

func (c *common) orignalPath(path string) string {
	return path[len(c.root):]
}

func (c *common) KV(path string) (map[string]string, error) {
	panic("abstract method")
	return nil, nil
}

func (c *common) Watch(path string) (chan []Event, chan error) {
	panic("abstract method")
	return nil, nil
}

func (c *common) Close() error {
	panic("abstract method")
	return nil
}
