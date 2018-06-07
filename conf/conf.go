package main

import (
	"encoding/json"
	"fmt"
	// "io"
	"log"
	// "strings"
)

type Config struct {
	Store       string
	StoreConfig json.RawMessage
}

type MysqlConfig struct {
	addr     string
	db       string
	user     string
	password string
}

type RedisConfig struct {
	addr string
	db   string
}

var j = []byte(`[{"Store": "mysql","StoreConfig": {"addr": "127.0.0.1","db": "test","password": "admin","user": "root"},{"Store": "redis","StoreConfig": {"addr": "192.168.1.33", "db": "0"}}]`)

func main() {
	var conf []Config
	err := json.Unmarshal(j, &conf)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	for _, c := range conf {
		var dst interface{}
		switch c.Store {
		case "mysql":
			dst = new(MysqlConfig)
		case "redis":
			dst = new(RedisConfig)
		}
		err := json.Unmarshal(c.StoreConfig, dst)
		if err != nil {
			log.Fatalln("Error: ", err)
		}
		fmt.Println(c.Store, dst)
	}

}
