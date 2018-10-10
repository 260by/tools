package main

import (
	"fmt"
	// "time"
	"github.com/260by/tools/sys/disk"
	"github.com/260by/tools/sys/cpu"
	"github.com/260by/tools/sys/mem"
	"github.com/260by/tools/sys/load"
	"github.com/260by/tools/sys/net"
	psdisk "github.com/shirou/gopsutil/disk"
)

func main()  {
	cpuUsage := cpu.Usage()
	fmt.Printf("CPU Used Percent: %.2f%%\n", cpuUsage)
	diskUsage := disk.Usage()
	fmt.Println(diskUsage)
	memUsage := mem.Usage()
	fmt.Printf("Mem Used Percent: %.2f%%\n", memUsage)
	loadAvg := load.Avg()
	fmt.Printf("Load1: %v Load5: %v Load15: %v\n", loadAvg[0], loadAvg[1], loadAvg[2])
	tcpState := net.TCPState()
	fmt.Printf("TCP State: %v\n", tcpState)

	ioStats, err := psdisk.IOCounters("/dev/sdb1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("IO: %v\n", ioStats)
}