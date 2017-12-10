package ini

import (
	"github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIniToJSON(t *testing.T) {
	kv := map[string]string{
		"key": "val",
	}
	data, err := IniToJSON(kv)
	assert.Nil(t, err)
	assert.Equal(t, `{"key":"val"}`, string(data))
}

func TestArray(t *testing.T) {
	kv := map[string]string{
		"user.0.name":       "yaoguai",
		"user.0.age":        "25",
		"user.0.interest.0": "coding",
		"user.0.interest.1": "reading",
	}
	data, err := IniToJSON(kv)
	assert.Nil(t, err)
	v := jsoniter.Get(data, "user", 0, "name").ToString()
	assert.Equal(t, "yaoguai", v)
	v = jsoniter.Get(data, "user", 0, "age").ToString()
	assert.Equal(t, "25", v)
	v = jsoniter.Get(data, "user", 0, "interest", 0).ToString()
	assert.Equal(t, "coding", v)
	v = jsoniter.Get(data, "user", 0, "interest", 1).ToString()
	assert.Equal(t, "reading", v)
}

func TestMistakeArray(t *testing.T) {
	kv := map[string]string{
		"user.2.name":       "yaoguai",
		"user.2.age":        "25",
		"user.2.interest.0": "coding",
		"user.2.interest.1": "reading",
	}
	data, err := IniToJSON(kv)
	assert.Nil(t, err)
	v := jsoniter.Get(data, "user", 0, "name").ToString()
	assert.NotEqual(t, "yaoguai", v)
	v = jsoniter.Get(data, "user", 0, "age").ToString()
	assert.NotEqual(t, "25", v)
	v = jsoniter.Get(data, "user", "2", "interest", 0).ToString()
	assert.Equal(t, "coding", v)
	v = jsoniter.Get(data, "user", "2", "interest", 1).ToString()
	assert.Equal(t, "reading", v)
}

func BenchmarkIniToJson(b *testing.B) {
	kv := map[string]string{
		"user.2.name":       "yaoguai",
		"user.2.age":        "25",
		"user.2.interest.0": "coding",
		"user.2.interest.1": "reading",
	}
	for i := 0; i < b.N; i++ {
		IniToJSON(kv)
	}
}
