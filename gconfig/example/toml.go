package main

import (
	"fmt"
	"github.com/260by/tools/gconfig"
)

type Config struct {
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

func main() {
	var conf Config
	if err := gconfig.Parse("config.toml", &conf); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("ENV: %s\n", conf.Env)
	fmt.Printf("Service: %s:%v\n", conf.ServiceIP, conf.ServicePort)
	fmt.Printf("Database: %s:%v\n", conf.Database.Host, conf.Database.Port)
}
