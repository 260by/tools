package disk

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"syscall"
)

// Usage 获取磁盘使用率
func Usage() map[string]float64 {
	diskUsage := make(map[string]float64)
	mountPoints := getMounts()

	for _, mount := range mountPoints {
		usedPercent := getMountPercent(mount)
		diskUsage[mount] = usedPercent
	}
	return diskUsage
}

func getMountPercent(path string) (usedPercent float64) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		fmt.Println(err)
	}
	total := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	used := total - free
	return float64(used) / float64(total) * 100
}

func getMounts() (mounts []string) {
	contents, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 2 && (fields[2] == "ext3" || fields[2] == "ext4" || fields[2] == "xfs") {
			reBoot, _ := regexp.Match("/boot.*", []byte(fields[1]))
			reDocker, _ := regexp.Match(".*docker.*", []byte(fields[1]))
			if !reBoot && !reDocker {
				mounts = append(mounts, fields[1])
			}
		}
	}
	return
}
