package sober

import (
	"github.com/yaoguais/sober/ds"
)

var (
	DataSource ds.DataSource
)

func Get(key string) (string, error) {
	return DataSource.Get(key)
}

func String(key string) string {
	val, _ := DataSource.Get(key)
	return val
}