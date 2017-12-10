package sober

import (
	"github.com/stretchr/testify/assert"
	"github.com/yaoguais/sober/ds"
	"testing"
)

func init() {
	DataSource, _ = ds.NewFile("./ds/.kv.ini")
}

func TestSober(t *testing.T) {
	assert.Equal(t, "yaoguai", String("user.name"))
	assert.Equal(t, true, Bool("key"))
	const age = 25
	var v0 int = age
	assert.Equal(t, v0, Int("user.age"))
	var v1 int32 = age
	assert.Equal(t, v1, Int32("user.age"))
	var v2 int64 = age
	assert.Equal(t, v2, Int64("user.age"))
	var v3 uint = age
	assert.Equal(t, v3, Uint("user.age"))
	var v4 uint32 = age
	assert.Equal(t, v4, Uint32("user.age"))
	var v5 uint64 = age
	assert.Equal(t, v5, Uint64("user.age"))
	var v6 float32 = 0.000001
	assert.Equal(t, true, v6-Float32("float") == 0)
	var v7 float64 = 0.000001
	assert.Equal(t, true, v7-Float64("float") == 0)
	assert.Nil(t, LastError())
}

func BenchmarkSober(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get("user.interest.0")
	}
}
