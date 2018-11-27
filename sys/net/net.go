package net

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func loadData() []string {
	var str []string

	tcpFile := "/proc/net/tcp"

	fin, err := os.Open(tcpFile)

	defer fin.Close()

	if err != nil {
		fmt.Println(tcpFile, err)
		return str
	}

	r := bufio.NewReader(fin)

	for {
		buf, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		str = append(str, buf)
	}

	if len(str) > 0 {
		return str[1:]
	}
	return str
}

func hex2dec(hexstr string) string {
	i, _ := strconv.ParseInt(hexstr, 16, 0)
	return strconv.FormatInt(i, 10)
}

func hexToIP(hexstr string) (string, string) {
	var ip string
	if len(hexstr) != 8 {
		err := "parse error"
		return ip, err
	}

	i1, _ := strconv.ParseInt(hexstr[6:8], 16, 0)
	i2, _ := strconv.ParseInt(hexstr[4:6], 16, 0)
	i3, _ := strconv.ParseInt(hexstr[2:4], 16, 0)
	i4, _ := strconv.ParseInt(hexstr[0:2], 16, 0)
	ip = fmt.Sprintf("%d.%d.%d.%d", i1, i2, i3, i4)

	return ip, ""
}

func convertToIPPort(str string) (string, string) {
	l := strings.Split(str, ":")
	if len(l) != 2 {
		return str, ""
	}

	ip, err := hexToIP(l[0])
	if err != "" {
		return str, ""
	}

	return ip, hex2dec(l[1])
}

func removeAllSpace(l []string) []string {
	var ll []string
	for _, v := range l {
		if v != "" {
			ll = append(ll, v)
		}
	}

	return ll
}

var tcpStatuses = map[string]string{
	"01": "ESTABLISHED",
	"02": "SYN_SENT",
	"03": "SYN_RECV",
	"04": "FIN_WAIT1",
	"05": "FIN_WAIT2",
	"06": "TIME_WAIT",
	"07": "CLOSE",
	"08": "CLOSE_WAIT",
	"09": "LAST_ACK",
	"0A": "LISTEN",
	"0B": "CLOSING",
}

// TCPState 统计TCP ESTABLISHED TIME_WAIT数量
func TCPState() map[string]int {
	tcpStateNum := make(map[string]int)
	lines := loadData()
	var established, timeWait int

	for _, line := range lines {
		l := removeAllSpace(strings.Split(line, " "))
		if tcpStatuses[l[3]] == "ESTABLISHED" {
			established++
		}
		if tcpStatuses[l[3]] == "TIME_WAIT" {
			timeWait++
		}
	}
	// fmt.Printf("ESTABLISHED: %v\nTIME_WAIT: %v\n", established, timeWait)
	tcpStateNum["ESTABLISHED"] = established
	tcpStateNum["TIME_WAIT"] = timeWait
	return tcpStateNum
}
