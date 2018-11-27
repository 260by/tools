package ini

import (
	"gopkg.in/gcfg.v1"
)

func ParseINI(f []byte, config interface{}) error {
	return gcfg.ReadStringInto(config, string(f))
}
