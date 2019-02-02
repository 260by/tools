package main

import (
	"net"
	"sync"
	"log"
	"io/ioutil"
	"strings"
	"fmt"
	"time"
	"flag"
	"net/http"
	"net/url"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"
)

func main()  {
	file := flag.String("f", "", "Proxy addr filename")
	flag.Parse()

	f, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalln(err)
	}
	ips := strings.Split(strings.TrimSuffix(string(f), "\n"), "\n")

	proxy := checkProxy(ips)

	for _, v := range proxy {
		if accessCheck(v, "http://www.vogued.cn") == 200 {
			fmt.Println(v)
		}
	}
}

func accessCheck(proxy, accessURL string) int {
	client := newHTTPClient(proxy)
	req, err := http.NewRequest("GET", accessURL, nil)
	if err != nil {
		// log.Println(err)
		return 503
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		// log.Println(err)
		return 503
	}

	return resp.StatusCode
}

func ping(addr string) bool {
	c, err := net.DialTimeout("tcp", addr, time.Millisecond*500)
	if err != nil {
		// log.Printf("Check address %v is error, %v", addr, err)
		return false
	}
	defer c.Close()
	return true
}

func checkProxy(proxy []string) []string {
	var proxyList []string
	var wg sync.WaitGroup
	for _, v := range proxy {
		
		wg.Add(1)
		go func(addr string) {
			defer wg.Add(-1)

			result := ping(addr)
			if result {
				// s := fmt.Sprintf("%s://%s", mold, addr)
				proxyList = append(proxyList, addr)
			}
		}(v)
	}
	wg.Wait()
	return proxyList
}

func newHTTPClient(proxyAddr string) *http.Client {
	proxyAddr = "http://" + proxyAddr
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return nil
	}

	netTransport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Millisecond*1000)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,                      //每个host最大空闲连接
		ResponseHeaderTimeout: time.Millisecond * 2000, //数据收发超时
	}

	return &http.Client{
		Timeout:   time.Second * 2,
		Transport: netTransport,
	}
}