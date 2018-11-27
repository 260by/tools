package main

import (
	"fmt"
	"github.com/260by/tools/pssh"
)

func main() {
	ssh := &pssh.Server{
		Addr:    "192.168.1.173",
		Port:    "22",
		User:    "root",
		KeyFile: "/home/user/.ssh/id_rsa",
	}

	stdout, err := ssh.Command("sudo /sbin/ifconfig")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(strings.Split(stdout, " "))
	fmt.Println(stdout)
}
