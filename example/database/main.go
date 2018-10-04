package main

import (
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"time"
)

type User struct {
	Id       int       `xorm:"not null pk autoincr INT(11)"`
	Username string    `xorm:"not null VARCHAR(32)"`
	Birthday time.Time `xorm: "DATE"`
	Sex      string    `xorm:"CHAR(1)"`
	Address  string    `xorm:"VARCHAR(256)"`
}

func main() {
	engine, err := xorm.NewEngine("mysql", "root:power123@tcp(192.168.1.251:3306)/ttt?charset=utf8mb4")
	if err != nil {
		log.Fatalln(err)
	}

	if err := engine.Ping(); err != nil {
		log.Fatalln(err)
	}

	//日志打印SQL
	engine.ShowSQL(true)

	result, err := engine.IsTableExist(&User{})
	if err != nil {
		log.Fatalln(err)
	}

	if !result {
		err := engine.CreateTables(&User{})
		if err != nil {
			log.Fatalln(err)
		}
	}
}


