package sober

import (
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	soberds "github.com/yaoguais/sober/ds"
	"testing"
)

func init() {
	ds, _ = soberds.NewFile("./ds/.kv.ini")
}

func TestSober(t *testing.T) {
	assert.Equal(t, "yaoguai", String("user.name"))
	assert.Equal(t, false, cast.ToBool(String("key")))
	const age = 25
	var v0 int = age
	assert.Equal(t, v0, cast.ToInt(String("user.age")))
	var v1 int32 = age
	assert.Equal(t, v1, cast.ToInt32(String("user.age")))
	var v2 int64 = age
	assert.Equal(t, v2, cast.ToInt64(String("user.age")))
	var v3 uint = age
	assert.Equal(t, v3, cast.ToUint(String("user.age")))
	var v4 uint32 = age
	assert.Equal(t, v4, cast.ToUint32(String("user.age")))
	var v5 uint64 = age
	assert.Equal(t, v5, cast.ToUint64(String("user.age")))
	var v6 float32 = 0.000001
	assert.Equal(t, true, v6-cast.ToFloat32(String("float")) == 0)
	var v7 float64 = 0.000001
	assert.Equal(t, true, v7-cast.ToFloat64(String("float")) == 0)
}

func BenchmarkSober(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get("user.interest.0")
	}
}
