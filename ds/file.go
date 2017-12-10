package ds

import (
	"errors"
	"github.com/djherbis/times"
	"github.com/go-ini/ini"
	"github.com/yaoguais/sober"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	validKey *regexp.Regexp
)

type Map map[string]interface{}
type Array []interface{}

func init() {
	validKey, _ = regexp.Compile("[a-zA-Z](\\.?[0-9a-zA-Z_])*")
}

type File struct {
	name  string
	ctime time.Time
	cfg   *ini.File
	c     chan struct{}
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
		c:     make(chan struct{}),
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
	f.RLock()
	defer f.RUnlock()

	if !validKey.Match([]byte(key)) {
		return "", errors.New("illegal key")
	}

	section := ""
	if i := strings.LastIndex(key, "."); i > -1 {
		section = key[:i]
	}

	if v, err := f.cfg.Section(section).GetKey(key); err != nil {
		return "", err
	} else {
		return v.Value(), nil
	}
}

func (f *File) Set(key, val string) error {
	return errors.New("forbid set")
}

func (f *File) JSON() ([]byte, error) {
	f.RLock()
	defer f.RUnlock()

	kv := make(map[string]string)
	for _, section := range f.cfg.Sections() {
		keys := section.Keys()
		for _, v := range keys {
			kv[v.Name()] = v.Value()
		}
	}

	return sober.IniToPrettyJSON(kv)
}

func (f *File) Watch() (chan struct{}, chan error) {
	c := make(chan struct{}, 1)
	e := make(chan error, 1)
	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-f.c:
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

func (f *File) Close() error {
	if f.c == nil {
		return errors.New("close twice")
	}

	close(f.c)
	f.c = nil

	return nil
}
