package metrics

import (
	"sentinel/internal/models"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)


func Collect(serverId string)(models.Metrics, error){
	hostInfo, _ := host.Info()
	cpuUsage, _ := cpu.Percent(1*time.Second, false)
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")

	// data to be sent
	data := models.Metrics{
		ServerId: serverId,
		Hostname: hostInfo.Hostname,
		CPU: cpuUsage[0],
		RAM: memInfo.UsedPercent,
		Disk: diskInfo.UsedPercent,
		Uptime: hostInfo.Uptime,
		TimeStamp: time.Now().Unix(),

	}
	return  data, nil
}