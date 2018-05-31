package yaml_config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

/* server.yaml file
env: test
server:
  ip: 0.0.0.0
  port: 5000
database:
  driver: mysql
  dsn: 'username:password@tcp(192.168.1.86:3306)/db_name'
redis:
  addr: 192.168.1.153:6385
  db: 1
*/

/* Example
package main

import (
	"fmt"
	config "github.com/260by/tools/yaml_config"
)

func main() {
	configFile := "server.yaml"
	if err := config.LocadFile(configFile); err != nil {
		fmt.Println(err)
	}
	// fmt.Println(config.Config.Server.IP)
	fmt.Printf("%s:%d\n", config.Config.Server.IP, config.Config.Server.Port)
}
*/

type config struct {
	Env    string
	Server struct {
		IP   string
		Port int
	}
	Database struct {
		Driver string
		DSN    string
	}
}

var Config = config{}

func Load(data []byte) error {
	return yaml.Unmarshal(data, &Config)
}

func LocadFile(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return Load(data)
}
