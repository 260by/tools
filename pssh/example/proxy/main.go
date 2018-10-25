package main

import (
    "fmt"
    "github.com/260by/tools/pssh"
)

func main()  {
    ssh := &pssh.Server{
	 Addr:  "10.111.1.12",
	 Port:    "22",
     User:    "root",
     KeyFile: "/home/user/.ssh/id_rsa",
     Proxy: pssh.ProxyServer{
		 Addr:  "123.57.80.54",
		 Port:    "22",
         User:    "bot",
         KeyFile: "/home/keith/id_rsa",
     },
    }

    stdout, err := ssh.Command("ifconfig")
    if err != nil {
     fmt.Println(err)
    }
    fmt.Println(stdout)

    f, e := ssh.Get("/root", "tmp")
    if e != nil {
     fmt.Println(e)
    }
    if f {
     fmt.Println("OK")
    }

    f, e = ssh.Put("tmp/a.txt", "/root")
    if e != nil {
        fmt.Println(e)
    }
    if f {
        fmt.Println("OK")
    }

}