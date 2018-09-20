package toml

import (
	"bytes"
	"github.com/BurntSushi/toml"
)

func ParseTOML(f []byte, config interface{}) error {
	if _, err := toml.DecodeReader(bytes.NewReader(f), config); err != nil {
		return err
	}
	return nil
}