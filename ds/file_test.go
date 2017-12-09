package ds

import (
	"github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestToJson(t *testing.T) {
	args := Args{
		DS: "file://.kv.ini",
	}
	ds, err := Provider(args)

	// {"key":"val","user":{"name":"yaoguai","age":"25","interest":["coding","reading"]},"database":{"master":{"host":"127.0.0.1","port":"3306"}},"parent":{"name":"rick","son":{"name":"carl"},"daughter":{"name":"judith"}}}
	data, err := ds.JSON()
	assert.Nil(t, err)
	v := jsoniter.Get(data, "user", "interest", 0).ToString()
	assert.Equal(t, "coding", v)
	v = jsoniter.Get(data, "database", "master", "host").ToString()
	assert.Equal(t, "127.0.0.1", v)
}
