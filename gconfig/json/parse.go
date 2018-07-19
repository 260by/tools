package json

import (
	"bytes"
	"encoding/json"
)

func ParseJSON(f []byte, config interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(f))
	decoder.UseNumber()
	return decoder.Decode(config)
}