package ds

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/djherbis/times"
	"github.com/go-ini/ini"
)

type File struct {
	name  string
	ctime time.Time
	cfg   *ini.File
	done  chan struct{}
	sync.RWMutex
}

func NewFile(name string) (*File, error) {
	cfg, err := ini.Load(name)
	if err != nil {
		return nil, err
	}

	t, err := ctime(name)
	if err != nil {
		return nil, err
	}

	f := &File{
		name:  name,
		cfg:   cfg,
		ctime: t,
		done:  make(chan struct{}),
	}

	return f, nil
}

func ctime(name string) (time.Time, error) {
	t, err := times.Stat(name)
	if err != nil {
		return time.Time{}, err
	}
	return t.ChangeTime(), nil
}

func (f *File) Get(key string) (string, error) {
	if !validKey.Match([]byte(key)) {
		return "", ErrIllegalKey
	}

	section := ""
	if i := strings.LastIndex(key, "."); i > -1 {
		section = key[:i]
	}

	f.RLock()
	defer f.RUnlock()

	if v, err := f.cfg.Section(section).GetKey(key); err != nil {
		return "", err
	} else {
		return v.Value(), nil
	}
}

func (f *File) Set(key, val string) error {
	return errors.New("forbid set")
}

func (f *File) Value() (string, error) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	f.RLock()
	_, err := f.cfg.WriteTo(w)
	f.RUnlock()

	return b.String(), err
}

func (f *File) Watch() (chan struct{}, chan error) {
	c := make(chan struct{}, 1)
	e := make(chan error, 1)

	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-f.done:
				return
			case <-t.C:
				v, err := ctime(f.name)
				if err != nil {
					e <- err
				} else if v != f.ctime {
					f.Lock()
					if err := f.cfg.Reload(); err != nil {
						e <- err
					} else {
						f.ctime = v
						c <- struct{}{}
					}
					f.Unlock()
				}
			}
		}
	}()

	return c, e
}

func (f *File) Feedback(error bool, message string) error {
	return nil
}

func (f *File) Close() error {
	if f.done == nil {
		return errors.New("close twice")
	}

	close(f.done)
	f.done = nil

	return nil
}
