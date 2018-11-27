## 邮件发送工具
发送普通文本或还附件邮件

### Install
    go get -u github.com/260by/tools/mail

### Quick Start
    package main

    import (
        "github.com/260by/tools/mail"
        "fmt"
    )

    func main()  {
        mail := mail.Server{
            Addr: "smtp.qq.com",
            Port: 25,
            User: "admin@qq.com",
            Password: "password",
        }

        // e := mail.Send("title", "test mail", []string{"zengming@juwang.cn"})
        // if e != nil {
        // 	panic(e)
        // } else {
        // 	fmt.Println("Send OK")
        // }

        err := mail.SendAttach("attach test", "attach test mail", []string{"r.go"}, []string{"zengming@juwang.cn"})
        if err != nil {
            panic(err)
        } else {
            fmt.Println("Send OK")
        }
    }