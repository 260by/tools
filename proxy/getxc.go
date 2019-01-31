package main

import (
	// "os"
	"errors"
	"fmt"
	// "math/rand"
	"github.com/PuerkitoBio/goquery"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
	"flag"
	"io/ioutil"
	"strings"
	"log"
)

const (
	baseURL = "http://www.xicidaili.com/nn/"
)

type Proxy struct {
	IP   string
	Port string
	Mold string
}

func main() {
	num := flag.Int("num", 5, "Get xicidaili page number.")
	flag.Parse()
	// t := time.Now()
	proxy, err := GetProxy(*num)
	if err != nil {
		log.Fatalln(err)
}
	proxyList := checkProxy(proxy)
	// for _, p := range proxyList {
	// 	fmt.Println(p)
	// }

	str := strings.Replace(strings.Trim(fmt.Sprint(proxyList), "[]"), " ", "\n", -1)
	filename := "/tmp/proxy.txt"
	err = ioutil.WriteFile(filename, []byte(str), 0644)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Printf("File save to %s", filename)
	}
}

// GetProxy 获取代理地址, count为获取的页数
func GetProxy(count int) ([]Proxy, error) {
	/*
		var proxy []Proxy
		var wg sync.WaitGroup
		for page := 1; page <= count; page++ {
			wg.Add(1)
			go func(page int)  {
				defer wg.Done()

				url := baseURL + strconv.Itoa(page)
				res, err := getResponse(url)
				if err != nil {
					fmt.Println(err)
					return
				}
				// fmt.Println(res.StatusCode)
				if res.StatusCode == 200 {
					dom, err := goquery.NewDocumentFromResponse(res)
					if err != nil {
						fmt.Println(err)
						return
					}

					dom.Find("#ip_list tr").Each(func(i int, context *goquery.Selection)  {
						resDom := context.Find("td")
						ip := resDom.Eq(resDom.Length()-9).Text()
						port := resDom.Eq(resDom.Length()-8).Text()
						mold := resDom.Eq(resDom.Length()-5).Text()
						// fmt.Printf("IP: %v\tPort: %v\t Type: %v\n", ip, port, mold)
						if mold == "HTTP" || mold == "HTTPS" {
							proxy = append(proxy, Proxy{IP: ip, Port: port, Mold: mold})
						}
					})
				}
			}(page)
		}
		wg.Wait()
		return proxy, nil
	*/
	// oneProxy, err := getOnePageProxy()
	// if err != nil {
	// 	return nil, err
	// }

	// println("Check")
	// proxyList := checkProxy(oneProxy)
	// fmt.Println(proxyList)

	var proxy []Proxy
	for page := 1; page <= count; page++ {
		url := baseURL + strconv.Itoa(page)

		// rand.Seed(time.Now().Unix() + int64(page))
		// random := rand.Intn(len(proxyList))

		// println(proxyList[random])
		// proxyAddr := proxyList[random]
		res, err := getResponse(url, "")
		if err != nil {
			return nil, err
		}

		// fmt.Println("URL: ", url)
		// fmt.Println(res.StatusCode)
		if res.StatusCode == 200 {
			dom, err := goquery.NewDocumentFromResponse(res)
			if err != nil {
				return nil, err
			}
			dom.Find("#ip_list tr").Each(func(i int, context *goquery.Selection) {
				resDom := context.Find("td")
				ip := resDom.Eq(resDom.Length() - 9).Text()
				port := resDom.Eq(resDom.Length() - 8).Text()
				mold := resDom.Eq(resDom.Length() - 5).Text()
				// fmt.Printf("IP: %v\tPort: %v\t Type: %v\n", ip, port, mold)
				if mold == "HTTP" || mold == "HTTPS" {
					proxy = append(proxy, Proxy{IP: ip, Port: port, Mold: mold})
				}
			})
		}
		time.Sleep(3000 * time.Millisecond)
	}

	return proxy, nil

}

func getOnePageProxy() ([]Proxy, error) {
	var proxy []Proxy

	url := baseURL + strconv.Itoa(1)
	res, err := getResponse(url, "")
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		dom, err := goquery.NewDocumentFromResponse(res)
		if err != nil {
			return nil, err
		}
		dom.Find("#ip_list tr").Each(func(i int, context *goquery.Selection) {
			resDom := context.Find("td")
			ip := resDom.Eq(resDom.Length() - 9).Text()
			port := resDom.Eq(resDom.Length() - 8).Text()
			mold := resDom.Eq(resDom.Length() - 5).Text()
			if port == "80" {
				return
			}

			if mold == "HTTP" || mold == "HTTPS" {
				proxy = append(proxy, Proxy{IP: ip, Port: port, Mold: mold})
			}
		})
	} else {
		err := errors.New("Request on page error" + res.Status)
		return nil, err
	}

	return proxy, nil
}

// func getResponse(url string) (*http.Response, error) {
// 	client := &http.Client{}
// 	request, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")
// 	return client.Do(request)
// }

func getResponse(url, proxyAddr string) (*http.Response, error) {
	var client *http.Client
	if proxyAddr != "" {
		client = newHTTPClient(proxyAddr)
	} else {
		client = &http.Client{}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")
	return client.Do(request)
}

func getPageCount() (int, error) {
	response, err := getResponse(baseURL, "")
	if err != nil {
		return 0, err
	}
	dom, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return 0, err
	}
	resDom := dom.Find(".pagination a")
	pageCount, _ := strconv.Atoi(resDom.Eq(resDom.Length() - 2).Text())
	return pageCount, nil
}

func ping(addr string) bool {
	// c, err := net.Dial("tcp", addr)
	c, err := net.DialTimeout("tcp", addr, time.Millisecond*500)
	if err != nil {
		return false
	}
	// c.SetDeadline(time.Now().Add(500 * time.Millisecond))
	defer c.Close()
	return true
}

func checkProxy(proxy []Proxy) []string {
	var proxyList []string
	var wg sync.WaitGroup
	for _, p := range proxy {
		addr := p.IP + ":" + p.Port
		wg.Add(1)
		go func(addr string) {
			defer wg.Add(-1)

			result := ping(addr)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			if result {
				proxyList = append(proxyList, addr)
			}
		}(addr)
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

	println(proxyAddr)
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
