package mem

import (
	// "fmt"
	"io/ioutil"
	"strconv"
	"strings"
	// "time"
)

// Usage 获取内存使用率
func Usage() (memUsage float64) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}

	var memTotal, memFree, memBuffers, memCached int
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 2 {
			val, _ := strconv.Atoi(fields[1])
			switch fields[0] {
			case "MemTotal:":
				memTotal = val
			case "MemFree:":
				memFree = val
			case "Buffers:":
				memBuffers = val
			case "Cached:":
				memCached = val
			}
		}
	}
	memUsage = float64(memTotal-memFree-memBuffers-memCached) / float64(memTotal) * 100
	return
}
