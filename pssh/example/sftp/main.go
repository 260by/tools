package main

import (
	"fmt"
	"github.com/260by/tools/pssh"
)

func main() {
	ssh := &pssh.Server{
		Addr:    "192.168.1.118",
		Port:    "22",
		User:    "root",
		KeyFile: "/home/user/.ssh/id_rsa",
	}

	// 上传文件
	f, e := ssh.Put("tmp/20180609.tar.gz", "/root")
	if e != nil {
		fmt.Println(e)
	}
	if f {
		fmt.Println("OK")
	}

	// 下载文件
	f, e = ssh.Get("/data/logs/20180609.tar.gz", "tmp")
	if e != nil {
		fmt.Println(e)
	}
	if f {
		fmt.Println("OK")
	}
}
