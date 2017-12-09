package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/ds"
	puter "github.com/yaoguais/sober/output"
	"os"
	"os/signal"
	"syscall"
)

var (
	datasource string
	output     string
	token      string
	root       string
	debug      bool
)

func init() {
	flag.StringVar(&datasource, "datasource", "file://.env.ini", "data source, file://kv.ini, grpc://host:ip")
	flag.StringVar(&token, "token", "", "authorize token")
	flag.StringVar(&root, "root", "", "root path of keys, like /dev/blog/nginx/backend")
	flag.StringVar(&output, "output", "", "output, file://.config.json")
	flag.BoolVar(&debug, "debug", false, "enable debug")
}

func main() {
	flag.Parse()

	args := ds.Args{
		DS:    datasource,
		Token: token,
		Root:  root,
	}

	ds, err := ds.Provider(args)
	if err != nil {
		logrus.WithError(err).Error("database source provider")
		os.Exit(1)
	}

	data, err := ds.JSON()
	if err != nil {
		logrus.WithError(err).Error("convert to json")
		os.Exit(1)
	}

	o, err := puter.Provider(output)
	if err != nil {
		logrus.WithError(err).Error("output provider")
		os.Exit(1)
	}

	if err := o.Put(data); err != nil {
		logrus.WithError(err).Error("put failed")
		os.Exit(1)
	}

	c, errC := ds.Watch()

	go func() {
		for range c {
			logrus.Info("reload config")
			data, err := ds.JSON()
			if err != nil {
				logrus.WithError(err).Error("convert to json")
				continue
			}

			if err := o.Put(data); err != nil {
				logrus.WithError(err).Error("put failed")
			}
		}
	}()

	go func() {
		for err := range errC {
			logrus.WithError(err).Error("watch")
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logrus.WithField("signal", sig).Info("receive signal")
		ds.Close()
		done <- struct{}{}
	}()

	logrus.Info("awaiting signal")
	<-done
	logrus.Info("exit")
}
