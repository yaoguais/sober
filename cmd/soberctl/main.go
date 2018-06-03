package main

import (
	"flag"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/store"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	root  string
	rule  string
	etcd  string
	debug bool

	stor store.Store
)

func init() {
	flag.StringVar(&root, "root", "/config/center", "root for all keys")
	flag.StringVar(&rule, "rule", "^(\\/[a-zA-Z0-9_.-]+){4,}$", "key validate rule")
	flag.StringVar(&etcd, "etcd", "127.0.0.1:2379", "etcd addresse 127.0.0.1:2379,127.0.0.1:2381")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.Parse()
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
Examples:
Get config of a key
soberctl get key
soberctl get /dev/blog/nginx/backend

Set config of a key
soberctl set key [./${key}]
soberctl set /dev/blog/software/nginx
soberctl set /dev/blog/software/nginx /tmp/nginx.conf`)
	os.Exit(1)
}

func main() {
	initLog()
	initStore()

	args := flag.Args()
	if len(args) < 2 {
		Usage()
	}

	validKey, err := regexp.Compile(rule)
	if err != nil {
		logrus.WithError(err).Error("invalid rule")
		os.Exit(1)
	}

	key := args[1]
	if !validKey.Match([]byte(key)) {
		logrus.Errorf("key %s not match rule %s", key, rule)
		os.Exit(1)
	}

	switch args[0] {
	case "get":
		v, err := stor.Get(key)
		if err != nil {
			logrus.WithError(err).Error("read store")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Result:\n%s\n", v)
	case "set":
		if len(args) > 3 {
			Usage()
		}

		cfgFile := "." + key
		if len(args) == 3 {
			cfgFile = args[2]
		}
		bs, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			logrus.WithError(err).Error("read config file")
			os.Exit(1)
		}

		err = stor.Set(key, string(bs))
		if err != nil {
			logrus.WithError(err).Error("store failed")
			os.Exit(1)
		}

		fmt.Println("OK")
	default:
		Usage()
	}
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
