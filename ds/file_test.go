package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	args := Args{
		DS: "file://.kv.ini",
	}
	ds, err := Provider(args)
	assert.Nil(t, err)
	assert.NotNil(t, ds)

	tests := []struct {
		key string
		val string
	}{
		{"key", "val"},
		{"user.name", "yaoguai"},
		{"user.age", "25"},
		{"database.master.host", "127.0.0.1"},
		{"database.master.port", "3306"},
		{"parent.name", "rick"},
		{"parent.son.name", "carl"},
		{"parent.daughter.name", "judith"},
	}

	for _, c := range tests {
		v, err := ds.Get(c.key)
		assert.Nil(t, err)
		assert.Equal(t, c.val, v)
	}

	v, err := ds.Get("none")
	assert.NotNil(t, err)
	assert.Equal(t, "", v)
}
