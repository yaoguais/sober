package service

import (
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/authorize"
	"github.com/yaoguais/sober/dispatcher"
	"github.com/yaoguais/sober/kvpb"
	"github.com/yaoguais/sober/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"regexp"
	"strings"
)

var (
	ErrPermissionDenied = status.Error(codes.InvalidArgument, "permission denied")
)

type kv struct {
	stor     store.Store
	pathRule *regexp.Regexp
}

func NewKV(stor store.Store, pathRule *regexp.Regexp) *kv {
	return &kv{
		stor:     stor,
		pathRule: pathRule,
	}
}

func (o *kv) Get(ctx context.Context, req *kvpb.GetRequest) (*kvpb.GetResponse, error) {
	auth := authorize.Auth{
		Token: req.Token,
		Path:  req.Root,
	}
	if !authorize.Valid(auth) {
		return nil, ErrPermissionDenied
	}

	m, err := o.stor.KV(req.Root)
	if err != nil {
		return nil, err
	}

	return &kvpb.GetResponse{
		Kv: m,
	}, nil
}

func (o *kv) Watch(req *kvpb.WatchRequest, ws kvpb.KV_WatchServer) error {
	auth := authorize.Auth{
		Token: req.Token,
		Path:  req.Root,
	}
	if !authorize.Valid(auth) {
		return ErrPermissionDenied
	}

	c := dispatcher.NewClient(req.Root)
	logrus.WithField("client", c).Debug("watch")

	dispatcher.Register(c)
	defer dispatcher.UnRegister(c)
	defer func() {
		logrus.WithField("client", c).Debug("leave")
	}()

	for e := range c.Event() {
		evt := &kvpb.Event{
			Type: e.Type,
		}
		if err := ws.Send(evt); err != nil {
			logrus.WithField("client", c).WithError(err).Error("watch event")
			return err
		}

		logrus.WithField("client", c).WithField("event", e).Error("watch")
	}

	return nil
}
