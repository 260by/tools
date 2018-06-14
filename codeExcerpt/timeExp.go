package main

import (
	"time"
	"fmt"
)

func main()  {
	// fmt.Printf("当前时间戳: %v\n", time.Now().Unix())
	// fmt.Printf("当前格式化时间: %v\n", time.Now().Format("01"))
	// fmt.Printf("当前格式化时间: %v\n", time.Now().Format("20060102150405"))
	// fmt.Printf("时间戳转格式化时间: %v\n", time.Unix(1389058332, 0).Format("2006-01-02 15:04:05"))
	// fmt.Printf("一天前: %v\n", time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	// fmt.Printf("一小时前: %v\n", time.Now().Add((-1)*time.Hour).Format("200601021504"))
	t := time.Now()
	d, _ := time.ParseDuration("-24h")
	fmt.Printf("一天前: %v\n", t.Add(d).Format("2006-01-02 15:04:05"))
	fmt.Printf("一周前: %v\n", t.Add(d * 7))
	fmt.Printf("一月前: %v\n", t.Add(d * 30))
}