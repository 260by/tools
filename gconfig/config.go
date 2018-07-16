package gconfig

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ParseYamlFile(filename string, config interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(file, config)
}