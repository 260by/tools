package main

import (
	"fmt"
	"github.com/260by/tools/gssh"
)

func main() {
	ssh := &gssh.Server{
		Addr:    "192.168.1.173",
		Port:    "22",
		User:    "root",
		KeyFile: "/home/keith/public_key/local",
	}

	// 上传文件
	// f, e := ssh.Put("tmp/20180609.tar.gz", "/root")
	// if e != nil {
	// 	fmt.Println(e)
	// }
	// if f {
	// 	fmt.Println("OK")
	// }

	// 下载文件
	f, e := ssh.Get("/data/test-logs", "tmp")
	if e != nil {
		fmt.Println(e)
	}
	if f {
		fmt.Println("OK")
	}
}
