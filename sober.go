package sober

import (
	soberds "github.com/yaoguais/sober/ds"
)

var (
	ds soberds.DataSource
)

func SetDataSource(v soberds.DataSource) {
	ds = v
}

func Get(key string) (string, error) {
	return ds.Get(key)
}

func String(key string) string {
	val, _ := ds.Get(key)
	return val
}
