package ini

import (
	"github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestJSONToIni(t *testing.T) {
	bs, _ := ioutil.ReadFile("../ds/.kv.json")
	var v interface{}
	err := jsoniter.Unmarshal(bs, &v)
	assert.Nil(t, err)
	m, err := JSONToIni(v)
	assert.Nil(t, err)

	em := map[string]string{
		"key":                  "val",
		"float":                "0.000001",
		"user.age":             "25",
		"user.username":        "yaoguai",
		"user.interest.0":      "coding",
		"user.interest.1":      "reading",
		"database.master.port": "3306",
		"database.master.host": "127.0.0.1",
		"parent.name":          "rick",
		"parent.son.name":      "carl",
		"parent.daughter.name": "judith",
	}
	assert.Equal(t, em, m)
}

func TestSimpleObjectJSONToIni(t *testing.T) {
	bs := []byte(`
{
	"key":"val",
	"key1": "val1"
}
`)
	var v interface{}
	err := jsoniter.Unmarshal(bs, &v)
	assert.Nil(t, err)
	m, err := JSONToIni(v)
	assert.Nil(t, err)

	em := map[string]string{
		"key":  "val",
		"key1": "val1",
	}
	assert.Equal(t, em, m)
}

func TestSimpleArrayJSONToIni(t *testing.T) {
	bs := []byte(`
{
	"arr": ["val1", "val2"]
}
`)
	var v interface{}
	err := jsoniter.Unmarshal(bs, &v)
	assert.Nil(t, err)
	m, err := JSONToIni(v)
	assert.Nil(t, err)

	em := map[string]string{
		"arr.0": "val1",
		"arr.1": "val2",
	}
	assert.Equal(t, em, m)
}
