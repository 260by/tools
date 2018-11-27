package gconfig

import (
	"errors"
	"github.com/260by/tools/gconfig/ini"
	"github.com/260by/tools/gconfig/json"
	"github.com/260by/tools/gconfig/toml"
	"github.com/260by/tools/gconfig/yaml"
	"io/ioutil"
	"path"
	"strings"
)

func Parse(file string, config interface{}) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	fileSuffix := getFileSuffix(file)

	switch fileSuffix {
	case "json":
		err = json.ParseJSON(buf, config)
	case "yaml", "yml":
		err = yaml.ParseYaml(buf, config)
	case "ini":
		err = ini.ParseINI(buf, config)
	case "toml":
		err = toml.ParseTOML(buf, config)
	default:
		err = errors.New("Configration file format does not support")
	}
	return err
}

func getFileSuffix(f string) string {
	filePath := path.Base(f)
	fileSuffix := path.Ext(filePath)
	return strings.Trim(fileSuffix, ".")
}
