package client

import (
	"github.com/yaoguais/sober"
	"github.com/yaoguais/sober/kvpb"
	"golang.org/x/net/context"
	"io"
)

type Event kvpb.Event

type KV struct {
	token  string
	root   string
	kvc    kvpb.KVClient
	cancel bool
}

func NewKV(token, root string) *KV {
	return &KV{
		token:  token,
		root:   root,
		cancel: false,
	}
}

func (o *KV) SetKVC(kvc kvpb.KVClient) {
	o.kvc = kvc
}

func (o *KV) Get(path string) (map[string]string, error) {
	req := &kvpb.GetRequest{
		Token: o.token,
		Root:  o.root,
	}

	resp, err := o.kvc.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp.Kv, nil
}

func (o *KV) Watch(path string) (chan Event, chan error) {
	c := make(chan Event)
	errC := make(chan error)

	retry := sober.NewRetry(1, 60)

	go func() {
		req := &kvpb.WatchRequest{
			Token: o.token,
			Root:  o.root,
		}

	try:
		resp, err := o.kvc.Watch(context.Background(), req)
		if err != nil {
			errC <- err
		} else {
			for {
				evt, err := resp.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					errC <- err
					break
				} else {
					e := Event(*evt)
					c <- e

					retry.Reset()
				}
			}
		}

		if !o.cancel {
			retry.Wait()
			goto try
		}
	}()

	return c, errC
}

func (o *KV) Cancel() {
	o.cancel = true
}
