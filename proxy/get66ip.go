package main

import (
	"regexp"
	"net"
	"context"
	"fmt"
	"flag"
	"strings"
	"log"
	"time"
	"sync"
	"io/ioutil"
	"os"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

const (
	// baseURL = "http://www.66ip.cn/mo.php?sxb=&tqsl=%v&ports%%5B%%5D2=&ktip=&sxa=&radio=radio&submit=%%CC%%E1++%%C8%%A1"
	baseURL = "http://www.66ip.cn/nmtq.php?getnum=%v&isp=0&anonymoustype=0&start=&ports=&export=&ipaddress=&area=1&proxytype=0&api=66ip"
	userAgent = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"
)

func main() {
	num := flag.Int("num", 50, "Get ip number")
	flag.Parse()

	url := fmt.Sprintf(baseURL, *num)

	var res string
	chromeBrowser(url, &res)

	l := strings.Split(res, "\n")
	var proxys []string
	for _, v := range l {
		pattern := "([0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+:[0-9]+)"
		r := regexp.MustCompile(pattern)
		a := r.FindAllString(v, -1)
		if len(a) > 0 {
			proxys = append(proxys, a[0])
		}
	}

	newProxy := checkProxy(proxys)

	d := strings.Replace(strings.Trim(fmt.Sprint(newProxy), "[]"), " ", "\n", -1)
	filename := fmt.Sprintf("%s/proxy.txt", os.Getenv("HOME"))
	err := ioutil.WriteFile(filename, []byte(d), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func chromeBrowser(accessURL string, res *string) {
	ctxt, cancel := context.WithCancel(context.Background())
	// ctxt, cancel = context.WithTimeout(ctxt, 180 * time.Second)
	defer cancel()

	c, err := chromedp.New(ctxt, chromedp.WithRunnerOptions(
		runner.UserAgent(userAgent),
		// runner.Flag("disable-extensions-except", "/home/keith/alexaToolbar/4.0.3_0"),
		// runner.Flag("load-extension", "/home/keith/alexaToolbar/4.0.3_0"),
	))
	if err != nil {
		log.Println(err)
		return
	}

	// n := 3600
	err = c.Run(ctxt, chromedp.Tasks{
		chromedp.Navigate(accessURL),
		// chromedp.SendKeys(`//input[@name="tqsl"]`, strconv.Itoa(n)),
		chromedp.Sleep(3 * time.Second),
		// chromedp.Submit(`//input[@name="tqsl"]`),
		chromedp.WaitReady("//body"),
		
		chromedp.OuterHTML("//body", res),

	})
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Shutdown(ctxt)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Wait()
	if err != nil {
		log.Println(err)
		return
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