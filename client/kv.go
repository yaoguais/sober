package client

import (
	"io"
	"os"

	"github.com/yaoguais/sober/kvpb"
	soberetry "github.com/yaoguais/sober/retry"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var (
	hostname string
)

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}
}

type Event kvpb.Event

type KV struct {
	token  string
	key    string
	kvc    kvpb.KVClient
	cancel bool
}

func NewKV(token, key string, kvc kvpb.KVClient) *KV {
	return &KV{
		token: token,
		key:   key,
		kvc:   kvc,
	}
}

func (o *KV) Get() (string, error) {
	req := &kvpb.GetRequest{
		Key: o.key,
	}

	resp, err := o.kvc.Get(o.newContext(), req)
	if err != nil {
		return "", err
	}

	return resp.Value, nil
}

func (o *KV) Watch() (chan Event, chan error) {
	c := make(chan Event)
	errC := make(chan error)

	go func() {
		req := &kvpb.WatchRequest{
			Key: o.key,
		}

		retry := soberetry.New(1, 60)
	try:
		resp, err := o.kvc.Watch(o.newContext(), req)
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

func (o *KV) newContext() context.Context {
	md := metadata.Pairs(
		"token", o.token,
		"hostname", hostname,
	)
	return metadata.NewOutgoingContext(context.Background(), md)
}
