## zip归档工具

### Install
    go get -u github.com/260by/tools/zip

### Quick Start
1. 压缩为zip

        package main

        import (
            "github.com/260by/tools/zip"
        )

        func main() {
            err := zip.Compress("/data/logs/nginx", "/data/logs/nginx.zip")
            if err != nil {
                panic(err)
            }
        }

2. 解压zip

        package main

        import (
            "github.com/260by/tools/zip"
        )

        func main() {
            err := zip.Extract("/data/logs/nginx.zip", "/data/logs")
            if err != nil {
                panic(err)
            }
        }