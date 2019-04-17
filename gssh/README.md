## pssh golang的ssh client工具

支持通过ssh执行命令，上传、下载文件，并支持通过跳板执行私网地址服务器

### Installation
    go get -u github.com/260by/tools/pssh

### Quick start

1. ssh

        package main

        import (
            "fmt"
            "github.com/260by/tools/gssh"
        )

        func main()  {
            ssh := &gssh.Server{
                Addr:    "192.168.1.173",
                Port:    "22",
                User:    "root",
                KeyFile: "/home/user/.ssh/id_rsa",
            }

            stdout, err := ssh.Command("sudo /sbin/ifconfig")
            if err != nil {
                fmt.Println(err)
            }

            fmt.Println(stdout)
        }

2. sftp

        package main

        import (
            "fmt"
            "github.com/260by/tools/gssh"
        )

        func main()  {
            ssh := &pssh.Server{
                Addr:   "192.168.1.118",
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
            f, e = ssh.Get("/data/logs/nginx.tar.gz", "tmp")
            if e != nil {
                fmt.Println(e)
            }
            if f {
                fmt.Println("OK")
            }
        }

3. proxy

        package main

        import (
            "fmt"
            "github.com/260by/tools/gssh"
        )

        func main()  {
            ssh := &pssh.Server{
            // 后端服务器配置信息
            Addr:    "10.111.1.12",
            Port:    "22",
            User:    "root",
            KeyFile: "/home/user/.ssh/id_rsa",
            // 代理服务器(跳板机)配置信息
            Proxy: pssh.ProxyServer{
                Addr:    "139.22.99.108",
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