package load

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// Avg 获取系统当前1分钟，5分钟，15分钟平均负载
func Avg() (loadAvg []float64) {
	contents, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return
	}

	fields := strings.Fields(string(contents))
	for i := 0; i <= 2; i++ {
		load, _ := strconv.ParseFloat(fields[i], 64)
		loadAvg = append(loadAvg, load)
	}

	return
}
