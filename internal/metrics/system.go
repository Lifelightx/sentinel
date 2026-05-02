package metrics

import (
	
	"os"
	"sentinel/internal/models"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func Collect() (models.Metrics, error) {
	hostInfo, _ := host.Info()
	cpuUsage, _ := cpu.Percent(1*time.Second, false)
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")

	// data to be sent
	data := models.Metrics{

		Hostname:  os.Getenv("HOST_NAME"),
		CPU:       cpuUsage[0],
		IPv4:      getIP(),
		RAM:       memInfo.UsedPercent,
		Disk:      diskInfo.UsedPercent,
		Uptime:    hostInfo.Uptime,
		TimeStamp: time.Now().Unix(),
	}
	return data, nil
}

func getIP() string {
	ip := os.Getenv("HOST_IP")
	if ip != "" {
		return ip
	}
	return "unknown"
}
