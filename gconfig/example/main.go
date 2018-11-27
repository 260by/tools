package main

import (
	"fmt"
	"github.com/260by/tools/gconfig"
)

type config struct {
	Env         string
	ServiceIP   string
	ServicePort int
	LogLevel    string
	Database    struct {
		Driver   string
		User     string
		Password string
		Host     string
		Port     int
		DBName   string
		Charset  string
	}
	Passport struct {
		EndPoint  string
		SecretID  string
		SecretKey string
	}
}

var conf config

func main() {
	filename := "config.json"
	err := gconfig.Parse(filename, &conf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listen%s:%v\n", conf.ServiceIP, conf.ServicePort)

	fmt.Println(conf.Database.DBName)
	fmt.Println(conf.Database.Host, conf.Database.Port)
}
