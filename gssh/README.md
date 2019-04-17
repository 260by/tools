## gssh golang的ssh client工具

支持通过ssh执行命令，上传、下载文件，并支持通过跳板执行私网地址服务器

### Installation
    go get -u github.com/260by/tools/gssh

### Quick start

1. ssh

```
package main

import (
	"fmt"
	"github.com/260by/tools/gssh"
)

func main() {
	ssh := &gssh.Server{
		Options: gssh.ServerOptions{
			Addr: "192.168.1.173",
			Port: "22",
			User: "root",
			KeyFile: "/root/.ssh/id_rsa",
		},
	}

	stdout, err := ssh.Command("sudo /sbin/ifconfig")
	if err != nil {
		panic(err)
	}
	fmt.Println(stdout)
}
```

2. sftp 支持文件或目录通配符

```
package main

import (
	"fmt"
	"github.com/260by/tools/gssh"
)

func main() {
	ssh := &gssh.Server{
		Options: gssh.ServerOptions{
			Addr: "192.168.1.173",
			Port: "22",
			User: "root",
			KeyFile: "/root/.ssh/id_rsa",
		},
	}

	// 上传文件
	err := ssh.Put("tttt1111.txt", "/tmp")
	if err != nil {
	 	panic(err)
	}

	// 下载文件
	err := ssh.Get("/data/test-logs", "tmp")
	if err != nil {
		panic(err)
	}
}
```

3. proxy

```
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

	stdout, err := ssh.Command("ls /data/logs")
	if err != nil {
		panic(err)
	}
	fmt.Println(stdout)

	err = ssh.Get("/data/logs", "/tmp")
	if err != nil {
		panic(err)
	}

	err = ssh.Put("tmp/a.txt", "/data")
	if err != nil {
	    panic(err)
	}
}
```