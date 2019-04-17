package main

import (
	"fmt"
	"github.com/260by/tools/gssh"
)

func main() {
	ssh := &gssh.Server{
		Options: gssh.ServerOptions{
			Addr: "10.111.1.12",
			Port: "22",
			User: "root",
			KeyFile: "/root/.ssh/internal",
		},
		ProxyOptions: gssh.ServerOptions{
			Addr: "123.43.34.9",
			Port: "22",
			User: "root",
			KeyFile: "/root/.ssh/id_rsa",
		},
	}

	stdout, err := ssh.Command("ls -l /data/logs")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(stdout)

	err = ssh.Get("/root", "tmp")
	if err != nil {
		fmt.Println(err)
	}

	err = ssh.Put("tmp/a.txt", "/root")
	if err != nil {
		fmt.Println(err)
	}
}
