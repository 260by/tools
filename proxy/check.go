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
		fmt.Println(v)
	}
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