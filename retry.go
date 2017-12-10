package sober

import (
	"time"
)

type Retry struct {
	start   time.Duration
	current time.Duration
	max     time.Duration
	c       chan struct{}
}

func NewRetry(start, max int) *Retry {
	return &Retry{
		start:   time.Duration(start) * time.Second,
		current: time.Duration(start) * time.Second,
		max:     time.Duration(max) * time.Second,
	}
}

func (r *Retry) Wait() {
	if r.c == nil {
		r.c = make(chan struct{})
	}

	c := make(chan struct{})
	go func() {
		time.Sleep(r.current)
		c <- struct{}{}
	}()

	select {
	case <-c:
		break
	case <-r.c:
		break
	}

	r.current = r.current * 2
	if r.current > r.max {
		r.current = r.max
	}
}

func (r *Retry) Reset() {
	r.current = r.start
}

func (r *Retry) Cancel() {
	if r.c != nil {
		close(r.c)
		r.c = nil
	}
}
