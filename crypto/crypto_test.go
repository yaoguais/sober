package crypto

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrytDecrypt(t *testing.T) {
	tests := []struct {
		Plain   string
		Secret  string
		Encrypt string
	}{
		{
			Plain:   "yaoguai's me",
			Secret:  "1234567890abcdefghijklmnopqrstuv",
			Encrypt: "50126a8306d875cc8994e50f",
		},
	}

	for _, v := range tests {
		o, err := encrypt([]byte(v.Plain), []byte(v.Secret))
		assert.Nil(t, err)
		assert.Equal(t, v.Encrypt, fmt.Sprintf("%x", o))
		p, err := decrypt(o, []byte(v.Secret))
		assert.Nil(t, err)
		assert.Equal(t, v.Plain, string(p))
	}
}

func TestEncodeDecode(t *testing.T) {
	keys := []string{
		"password",
		"secretKey",
		"secret",
	}
	SetProtect(keys)

	tests := []struct {
		Key     string
		Val     string
		Encrypt string
	}{
		{
			Key:     "key",
			Val:     "val",
			Encrypt: "val",
		},
		{
			Key:     "password",
			Val:     "123456",
			Encrypt: magic + "184136d0468f",
		},
		{
			Key:     "password",
			Val:     magic + "184136d0468f",
			Encrypt: magic + "184136d0468f",
		},
	}

	SetSecret([]byte("1234567890abcdefghijklmnopqrstuv"))

	for _, v := range tests {
		e := Encode(v.Key, v.Val)
		assert.Equal(t, v.Encrypt, e)
		s := Decode(v.Key, e)
		if strings.HasPrefix(v.Val, magic) {
			assert.NotEqual(t, v.Val, s)
		} else {
			assert.Equal(t, s, v.Val)
		}
	}
}
