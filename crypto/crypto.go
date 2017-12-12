package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"strings"
)

var (
	// Size 16, 24 or 32 support
	Secret    []byte
	defaultIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	protect   map[string]struct{}
	none      = struct{}{}
	magic     = "\000\001\020"
	magicLen  = len(magic)
)

func Protect(l []string) {
	protect = make(map[string]struct{})
	for _, v := range l {
		protect[strings.ToLower(v)] = none
	}
}

func Encode(k, v string) string {
	if v == "" {
		return ""
	}

	if Secret == nil {
		return v
	}

	if strings.HasPrefix(v, magic) {
		return v
	}

	if _, ok := protect[strings.ToLower(k)]; !ok {
		return v
	}

	if t, err := encrypt([]byte(v), Secret); err != nil {
		return v
	} else {
		return magic + hex.EncodeToString(t)
	}
}

func Decode(k, v string) string {
	if len(v) <= magicLen {
		return v
	}

	if !strings.HasPrefix(v, magic) {
		return v
	}

	p := v[magicLen:]

	t, err := hex.DecodeString(p)
	if err != nil {
		return v
	}

	if c, err := decrypt(t, Secret); err != nil {
		return v
	} else {
		return string(c)
	}
}

func encrypt(s, secret []byte) ([]byte, error) {
	c, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	e := cipher.NewCFBEncrypter(c, defaultIV)
	r := make([]byte, len(s))
	e.XORKeyStream(r, s)

	return r, nil
}

func decrypt(s, secret []byte) ([]byte, error) {
	c, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	d := cipher.NewCFBDecrypter(c, defaultIV)
	p := make([]byte, len(s))
	d.XORKeyStream(p, s)

	return p, nil
}
