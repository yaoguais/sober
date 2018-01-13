package ini

import (
	"errors"
	"github.com/spf13/cast"
	"strings"
)

func JSONToIni(v interface{}) (map[string]string, error) {
	m, ok := v.(Map)
	if !ok {
		return nil, errors.New("v must be map[string]interface{}")
	}
	ini := make(map[string]string)
	for k, v := range m {
		if err := jsonToIni(v, []string{k}, ini); err != nil {
			return nil, err
		}
	}

	return ini, nil
}

func jsonToIni(v interface{}, keys []string, m map[string]string) error {
	switch v.(type) {
	case Map:
		for k, v := range v.(Map) {
			var ks []string
			ks = append(ks, keys...)
			ks = append(ks, k)
			if err := jsonToIni(v, ks, m); err != nil {
				return err
			}
		}
	case Array:
		for i, v := range v.(Array) {
			var ks []string
			ks = append(ks, keys...)
			ks = append(ks, cast.ToString(i))
			if err := jsonToIni(v, ks, m); err != nil {
				return err
			}
		}
	default:
		key := strings.Join(keys, ".")
		m[key] = cast.ToString(v)
	}
	return nil
}
