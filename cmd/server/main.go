package main

import (
	"flag"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/authorize"
	"github.com/yaoguais/sober/dispatcher"
	"github.com/yaoguais/sober/kvpb"
	"github.com/yaoguais/sober/service"
	"github.com/yaoguais/sober/store"
	"google.golang.org/grpc"
)

var (
	addr    string
	root    string
	authkey string
	rule    string
	etcd    string
	debug   bool

	stor store.Store
)

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0:3333", "server listen address, host:ip")
	flag.StringVar(&root, "root", "/config/center", "root for all keys")
	flag.StringVar(&authkey, "authkey", "/prod/infrastructure/service/sober.authorize.toml", "key for saving sober authroize config")
	flag.StringVar(&rule, "rule", "^(\\/[a-zA-Z0-9_.-]+){4,}$", "key validate rule")
	flag.StringVar(&etcd, "etcd", "127.0.0.1:2379", "etcd addresse 127.0.0.1:2379,127.0.0.1:2381")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.Parse()
}

func main() {
	initLog()
	initStore()
	initAuthorize()
	go initDispatcher()
	go initServer()

	watch()
}

func initLog() {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.SetOutput(os.Stdout)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func initStore() {
	opts := clientv3.Config{
		Endpoints:   strings.Split(etcd, ","),
		DialTimeout: 5 * time.Second,
	}

	kv, err := clientv3.New(opts)
	if err != nil {
		logrus.WithError(err).Error("create ectd")
		os.Exit(1)
	}

	s, err := store.NewEtcd(kv)
	if err != nil {
		logrus.WithError(err).Error("create etcd store")
		os.Exit(1)
	}

	s.SetRoot(root)

	stor = s
}

func initAuthorize() {
	err := authorize.Start(stor, authkey)
	if err != nil {
		logrus.WithError(err).Error("start authorize")
		os.Exit(1)
	}
}

func initDispatcher() {
	dispatcher.Dispatch(stor)
}

func initServer() {
	keyRule, err := regexp.Compile(rule)
	if err != nil {
		logrus.WithError(err).Debug("illegal rule")
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.WithError(err).WithField("addr", addr).Error("listen")
		os.Exit(1)
	}

	logrus.WithField("addr", addr).Info("listen")

	kvSvc := service.NewKV(stor, keyRule)

	grpcServer := grpc.NewServer()
	kvpb.RegisterKVServer(grpcServer, kvSvc)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.WithError(err).Error("serve")
		os.Exit(1)
	}
}

func watch() {
	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logrus.WithField("signal", sig).Info("receive signal")
		stor.Close()
		done <- struct{}{}
	}()

	<-done
	logrus.Info("exit")
}
