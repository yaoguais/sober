package main

import (
	"flag"
	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/authorize"
	"github.com/yaoguais/sober/store"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var (
	addr  string
	root  string
	path  string
	rule  string
	etcd  string
	debug bool
)

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0:3333", "server listen address, host:ip")
	flag.StringVar(&root, "root", "/config/center", "root for all paths")
	flag.StringVar(&path, "path", "/prod/infrastructure/service/sober", "path for saving data")
	flag.StringVar(&rule, "rule", "^(\\/[a-zA-Z0-9_-]+){4,}$", "root validate rule")
	flag.StringVar(&etcd, "etcd", "127.0.0.1:2379", "etcd addresse 127.0.0.1:2379,127.0.0.1:2381")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.Parse()
}

func main() {
	var stor store.Store

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

	rule, err := regexp.Compile(rule)
	if err != nil {
		logrus.WithError(err).Debug("illegal rule")
		os.Exit(1)
	}

	stor = s.SetRoot(root).SetRule(rule)

	err = authorize.Start(stor, path)
	if err != nil {
		logrus.WithError(err).Error("start authorize")
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logrus.WithField("signal", sig).Info("receive signal")
		stor.Close()
		done <- struct{}{}
	}()

	logrus.Info("awaiting signal")
	<-done
	logrus.Info("exit")
}
