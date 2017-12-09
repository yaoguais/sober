package main

import (
	"flag"
	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/store"
	"os"
	"os/signal"
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
}

func main() {
	etcdOpts := clientv3.Config{
		Endpoints:   strings.Split(etcd, ","),
		DialTimeout: 5 * time.Second,
	}
	store, err := store.NewEtcd(root, rule, path, etcdOpts)
	if err != nil {
		logrus.WithError(err).Error("create store")
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logrus.WithField("signal", sig).Info("receive signal")
		store.Close()
		done <- struct{}{}
	}()

	logrus.Info("awaiting signal")
	<-done
	logrus.Info("exit")
}
