package sober

import (
	"github.com/yaoguais/sober/ds"
	"strconv"
)

var (
	DataSource ds.DataSource
	lastErr    error
)

func Get(key string) (string, error) {
	return DataSource.Get(key)
}

func String(key string) string {
	v, err := DataSource.Get(key)
	if err != nil {
		lastErr = err
	}
	return v
}

func Bool(key string) bool {
	v, err := DataSource.Get(key)
	if err != nil {
		lastErr = err
	}
	return v != "" && v != "0"
}

func Int(key string) int {
	v, err := DataSource.Get(key)
	if err != nil {
		lastErr = err
	}
	n, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		lastErr = err
	}
	return int(n)
}

func Int32(key string) int32 {
	v, err := DataSource.Get(key)
	if err != nil {
		lastErr = err
	}
	n, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		lastErr = err
	}
	return int32(n)
}

func Int64(key string) int64 {
	v, err := DataSource.Get(key)
	if err != nil {
		lastErr = err
	}
	n, err := strconv.ParseInt(v, 10, 64)
	return n
}

func Uint(key string) uint {
	v, err := DataSource.Get(key)
	if err != nil {
		lastErr = err
	}
	n, err := strconv.ParseUint(v, 10, 0)
	if err != nil {
		lastErr = err
	}
	return uint(n)
}

func Uint32(key string) uint32 {
	v, err := DataSource.Get(key)
	if err != nil {
		lastErr = err
	}
	n, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		lastErr = err
	}
	return uint32(n)
}

func Uint64(key string) uint64 {
	v, err := DataSource.Get(key)
	if err != nil {
		lastErr = err
	}
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		lastErr = err
	}
	return n
}

func Float32(key string) float32 {
	v, err := DataSource.Get(key)
	if err != nil {
		lastErr = err
	}
	n, err := strconv.ParseFloat(v, 32)
	if err != nil {
		lastErr = err
	}
	return float32(n)
}

func Float64(key string) float64 {
	v, err := DataSource.Get(key)
	if err != nil {
		lastErr = err
	}
	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		lastErr = err
	}
	return n
}

func LastError() error {
	err := lastErr
	if lastErr != nil {
		lastErr = nil
	}
	return err
}
