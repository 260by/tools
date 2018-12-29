package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	// "time"
)

func main() {
	proxyFile := "ip.txt"
	f, err := ioutil.ReadFile(proxyFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	proxy := strings.Split(string(f), "\n")

	baseURL := "http://www.vogued.cn"

	for _, p := range proxy {
		p = "http://" + p
		proxyAddr, _ := url.Parse(p)
		fmt.Println("Proxy: ", p)

		netTransport := &http.Transport{
			Proxy: http.ProxyURL(proxyAddr),
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Millisecond*1000)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   10,                      //每个host最大空闲连接
			ResponseHeaderTimeout: time.Millisecond * 3000, //数据收发超时
		}

		client := &http.Client{
			Transport: netTransport,
			Timeout:   time.Second * 5,
		}

		request, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")

		resp, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(resp.StatusCode)
		time.Sleep(10 * time.Second)
	}
}
