package gconfig

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"
	"github.com/260by/tools/gconfig/yaml"
	"github.com/260by/tools/gconfig/json"
	"github.com/260by/tools/gconfig/ini"
)

func Parse(f string, config interface{}) error {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	fileSuffix := getFileSuffix(f)

	switch fileSuffix {
	case "json":
		err = json.ParseJSON(buf, config)
	case "yaml", "yml":
		err = yaml.ParseYaml(buf, config)
	case "ini":
		err = ini.ParseINI(buf, config)
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