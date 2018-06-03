package decode

import (
	"fmt"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
	"github.com/json-iterator/go"
)

func Decode(key, value string, v interface{}) error {
	switch e := path.Ext(key); e {
	case ".json":
		return jsoniter.UnmarshalFromString(value, v)
	case ".toml":
		_, err := toml.Decode(value, v)
		return err
	case ".yaml", ".yml":
		return yaml.Unmarshal([]byte(value), v)
	default:
		return fmt.Errorf("filetype '%s' not support", e)
	}
}
