package yaml

import (
	"gopkg.in/yaml.v2"
)

func ParseYaml(f []byte, config interface{}) error {
	return yaml.Unmarshal(f, config)
}