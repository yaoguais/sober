package main

import (
	"flag"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/yaoguais/sober/ds"
	"github.com/yaoguais/sober/ini"
	"github.com/yaoguais/sober/store"
	"io/ioutil"
	"os"
	gopath "path"
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
	flag.StringVar(&root, "root", "/config/center", "root for all paths")
	flag.StringVar(&rule, "rule", "^(\\/[a-zA-Z0-9_-]+){4,}$", "root validate rule")
	flag.StringVar(&etcd, "etcd", "127.0.0.1:2379", "etcd addresse 127.0.0.1:2379,127.0.0.1:2381")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.Parse()
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
Examples:
Get config of a path
soberctl get path
soberctl get /dev/blog/nginx/backend

Set config of a path
soberctl set path [./${path}.json]
soberctl set /dev/blog/nginx/backend
soberctl set /dev/blog/nginx/backend /tmp/config.json`)
	os.Exit(1)
}

func main() {
	initLog()
	initStore()

	args := flag.Args()
	if len(args) < 2 {
		Usage()
	}

	validPath, err := regexp.Compile(rule)
	if err != nil {
		logrus.WithError(err).Error("invalid rule")
		os.Exit(1)
	}

	path := args[1]
	if !validPath.Match([]byte(path)) {
		logrus.Errorf("path %s not match rule %s", path, rule)
		os.Exit(1)
	}

	switch args[0] {
	case "get":
		kv, err := stor.KV(path)
		if err != nil {
			logrus.WithError(err).Error("read store")
			os.Exit(1)
		}
		kv = ds.ReplaceToDotKey(kv)
		bs, err := ini.IniToPrettyJSON(kv)
		if err != nil {
			logrus.WithError(err).Error("convert to json failed")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Result:\n%s\n", string(bs))
	case "set":
		if len(args) > 3 {
			Usage()
		}

		cfgFile := gopath.Join(".", path, "config.json")
		if len(args) == 3 {
			cfgFile = args[2]
		}
		bs, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			logrus.WithError(err).Error("read config file")
			os.Exit(1)
		}

		var v interface{}
		err = jsoniter.Unmarshal(bs, &v)
		if err != nil {
			logrus.WithError(err).Error("parse config file")
			os.Exit(1)
		}

		kv, err := ini.JSONToIni(v)
		if err != nil {
			logrus.WithError(err).Error("convert to ini failed")
			os.Exit(1)
		}
		kv = ds.ReplaceToSlashKey(kv)

		err = stor.Set(path, kv)
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
