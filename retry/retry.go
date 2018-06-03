package retry

import (
	"time"
)

type retry struct {
	start   time.Duration
	current time.Duration
	max     time.Duration
	done    chan struct{}
}

func New(start, max int) *retry {
	return &retry{
		start:   time.Duration(start) * time.Second,
		current: time.Duration(start) * time.Second,
		max:     time.Duration(max) * time.Second,
	}
}

func (r *retry) Wait() {
	if r.done == nil {
		r.done = make(chan struct{})
	}

	c := make(chan struct{}, 1)
	go func() {
		time.Sleep(r.current)
		c <- struct{}{}
	}()

	select {
	case <-c:
		break
	case <-r.done:
		return
	}

	r.current = r.current * 2
	if r.current > r.max {
		r.current = r.max
	}
}

func (r *retry) Reset() {
	r.current = r.start
}

func (r *retry) Cancel() {
	if r.done != nil {
		close(r.done)
		r.done = nil
	}
}
