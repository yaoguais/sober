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
	key        string
	token      string
	debug      bool

	dso ds.DataSource
)

func init() {
	flag.StringVar(&datasource, "datasource", "file://.env.ini", "data source, file://kv.ini, grpc://host:ip")
	flag.StringVar(&key, "key", "", "key of config value, e.g. /dev/blog/backend/go.json")
	flag.StringVar(&token, "token", "", "authorize token")
	flag.StringVar(&output, "output", "", "output, file://.config.json")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.Parse()
}

func main() {
	initLog()
	initDataSource()
	initOutput()

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

func initDataSource() {
	args := ds.Args{
		DS:    datasource,
		Token: token,
		Key:   key,
	}

	ds, err := ds.Provider(args)
	if err != nil {
		logrus.WithError(err).Error("database source provider")
		os.Exit(1)
	}

	dso = ds
}

func initOutput() {
	o, err := puter.Provider(output)
	if err != nil {
		logrus.WithError(err).Error("output provider")
		os.Exit(1)
	}

	v, err := dso.Get(key)
	if err != nil {
		logrus.WithError(err).Error("get data")
		os.Exit(1)
	}

	if err := o.Put([]byte(v)); err != nil {
		logrus.WithError(err).Error("put failed")
		os.Exit(1)
	}

	logrus.WithField("data", v).Debug("put data")

	c, errC := dso.Watch()

	go func() {
		for range c {
			logrus.Info("reload config")
			v, err := dso.Get(key)
			if err != nil {
				logrus.WithError(err).Error("get data")
				continue
			}

			if err := o.Put([]byte(v)); err != nil {
				logrus.WithError(err).Error("put failed")
			}

			logrus.WithField("data", v).Debug("put data")
		}
	}()

	go func() {
		for err := range errC {
			logrus.WithError(err).Error("watch")
		}
	}()
}

func watch() {
	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logrus.WithField("signal", sig).Info("receive signal")
		dso.Close()
		done <- struct{}{}
	}()

	logrus.Info("awaiting signal")
	<-done
	logrus.Info("exit")
}
