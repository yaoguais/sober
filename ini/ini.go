package ini

import (
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"regexp"
	"strconv"
	"strings"
)

type Map map[string]interface{}
type Array []interface{}

func IniToJSON(kv map[string]string) ([]byte, error) {
	if m, err := iniToMap(kv); err != nil {
		return nil, err
	} else {
		return jsoniter.Marshal(m)
	}
}

func IniToPrettyJSON(kv map[string]string) ([]byte, error) {
	if m, err := iniToMap(kv); err != nil {
		return nil, err
	} else {
		return jsoniter.MarshalIndent(m, "", "    ")
	}
}

func iniToMap(kv map[string]string) (Map, error) {
	root := make(Map)
	for k, v := range kv {
		if err := addToJSON(root, kv, k, v); err != nil {
			return nil, err
		}
	}

	return root, nil
}

func addToJSON(root interface{}, kv map[string]string, key, val string) error {
	ks := strings.Split(key, ".")
	p := root
	i := 0
	for ; i < len(ks)-1; i++ {
		switch p.(type) {
		case Map:
			real := p.(Map)
			n, inited := real[ks[i]]
			if !inited {
				size, isArr := array(kv, strings.Join(ks[0:i+1], "."))
				if isArr {
					real[ks[i]] = make(Array, size, size)
				} else {
					real[ks[i]] = make(Map)
				}
				p = real[ks[i]]
			} else {
				p = n
			}
		case Array:
			real := p.(Array)
			j, err := strconv.ParseUint(ks[i], 10, 64)
			if err != nil {
				return err
			}
			if real[j] == nil {
				size, isArr := array(kv, strings.Join(ks[0:i+1], "."))
				if isArr {
					real[j] = make(Array, size, size)
				} else {
					real[j] = make(Map)
				}
			}
			p = real[j]
		default:
			return fmt.Errorf("illegal type %T 1", p)
		}
	}

	switch p.(type) {
	case Map:
		p.(Map)[ks[i]] = val
	case Array:
		idx, _ := strconv.ParseUint(ks[i], 10, 64)
		if int(idx) >= len(p.(Array)) {
			return errors.New("index out of range")
		}
		p.(Array)[idx] = val
	default:
		return fmt.Errorf("illegal type %T 2", p)
	}

	return nil
}

func array(kv map[string]string, section string) (int, bool) {
	s := fmt.Sprintf("%s\\.[0-9]+(\\.?[0-9a-zA-Z_])*", section)
	reg, _ := regexp.Compile(s)
	pl := len(section) + 1
	km := make(map[string]struct{})

	for k, _ := range kv {
		if reg.Match([]byte(k)) {
			if i := strings.Index(k[pl:], "."); i > -1 {
				km[k[pl:pl+i]] = struct{}{}
			} else {
				km[k[pl:]] = struct{}{}
			}
		}
	}

	size := len(km)
	if size > 0 {
		var i int64
		for ; i < int64(size); i++ {
			if _, ok := km[strconv.FormatInt(i, 10)]; !ok {
				return 0, false
			}
		}
	}

	return size, size > 0
}
