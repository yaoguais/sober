package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/authorize"
	"github.com/yaoguais/sober/dispatcher"
	"github.com/yaoguais/sober/kvpb"
	"github.com/yaoguais/sober/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type kv struct {
	stor    store.Store
	keyRule *regexp.Regexp
}

func NewKV(stor store.Store, keyRule *regexp.Regexp) *kv {
	return &kv{
		stor:    stor,
		keyRule: keyRule,
	}
}

func (o *kv) Get(ctx context.Context, req *kvpb.GetRequest) (*kvpb.GetResponse, error) {
	if err := permit(ctx, req.Key); err != nil {
		return nil, err
	}

	v, err := o.stor.Get(req.Key)
	if err != nil {
		return nil, err
	}

	return &kvpb.GetResponse{
		Value: v,
	}, nil
}
func (o *kv) Set(ctx context.Context, req *kvpb.SetRequest) (*kvpb.SetResponse, error) {
	return nil, errors.New("not permit")
}

func (o *kv) Watch(req *kvpb.WatchRequest, w kvpb.KV_WatchServer) error {
	if err := permit(w.Context(), req.Key); err != nil {
		return err
	}

	c := dispatcher.NewClient(req.Key)
	logrus.WithField("client", c).Debug("watch")

	dispatcher.Register(c)
	defer dispatcher.UnRegister(c)
	defer func() {
		logrus.WithField("client", c).Debug("leave")
	}()

	for e := range c.Event() {
		evt := &kvpb.Event{}
		if err := w.Send(evt); err != nil {
			logrus.WithField("client", c).WithError(err).Error("watch")
			return err
		}

		logrus.WithField("client", c).WithField("event", e).Debug("watch")
	}

	return nil
}

func (o *kv) Feedback(ctx context.Context, req *kvpb.FeedbackRequest) (*kvpb.FeedbackResponse, error) {
	hostname, err := getHostname(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.Error {
		logrus.WithFields(logrus.Fields{
			"hostname": hostname,
			"key":      req.Key,
			"message":  req.Message,
		}).Error("feedback")
	} else {
		logrus.WithFields(logrus.Fields{
			"hostname": hostname,
			"key":      req.Key,
			"message":  req.Message,
		}).Info("feedback")

	}
	return &kvpb.FeedbackResponse{}, nil
}

func getIncomingMetadataValue(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata %s not found", key)
	}
	if len(md) > 0 && len(md[key]) > 0 {
		if len(md[key][0]) > 0 {
			return md["token"][0], nil
		}
		return "", fmt.Errorf("empty %s found", key)
	}
	return "", fmt.Errorf("%s not found", key)

}

func getToken(ctx context.Context) (string, error) {
	return getIncomingMetadataValue(ctx, "token")
}

func getHostname(ctx context.Context) (string, error) {
	return getIncomingMetadataValue(ctx, "hostname")
}

func permit(ctx context.Context, key string) error {
	token, err := getToken(ctx)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	v := authorize.Auth{
		Token: token,
		Key:   key,
	}
	if !authorize.Valid(v) {
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	return nil
}
