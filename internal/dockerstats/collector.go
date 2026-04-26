package dockerstats

import (
	"context"
	"encoding/json"

	"sentinel/internal/models"
	"strings"

	"github.com/docker/docker/api/types/container"
)

type statsJSON struct{
	CPUStats struct{
		CPUUsage struct{
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemUsage uint64 `json:"system_cpu_sage"`
		OnlineCpus uint64 `json:"online_cpus"`
	} `json:"cpu_stats"`
	PreCpuStats struct{
		CPUUsage struct{
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemUsage uint64 `json:"system_cpu_usage"`
	} `json:"pre_cpu_stats"`
	MemoryStats struct{
		Usage uint64 `json:"usage"`
		Limit uint64 `json:"limit"`
	} `json:"memory_stats"`
}

func Collect() ([]models.ContainerInfo, error){
	cli, err :=NewClient()
	if err != nil{
		return  nil, err
	}
	list ,err := cli.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil{
		return nil, err
	}
	var results []models.ContainerInfo
	
	for _, c := range list{
		name := ""
		if len(c.Names)>0{
			name = strings.TrimPrefix(c.Names[0], "/")
		}
		statsResp, err := cli.ContainerStats(context.Background(), c.ID, false)
		if err != nil{
			continue
		}
		var stats statsJSON
		if err := json.NewDecoder(statsResp.Body).Decode(&stats); err != nil{
			statsResp.Body.Close()
			continue
		}
		statsResp.Body.Close()
		//debug
		// log.Println(stats)
		inspect, err := cli.ContainerInspect(context.Background(), c.ID)
		health := "none"
		if err != nil && inspect.State != nil && inspect.State.Health != nil{
			health = inspect.State.Health.Status
		}
		memMB := float64(stats.MemoryStats.Usage) / 1024 / 1024
		memPct := 0.0
		if stats.MemoryStats.Limit > 0{
			memPct = (float64(stats.MemoryStats.Usage)) / float64(stats.MemoryStats.Limit) * 100
		}
		id := c.ID
		if len(id)>12{
			id = id[:12]
		}
		results = append(results, models.ContainerInfo{
			ID: id,
			Name: name,
			State: c.State,
			Status: c.Status,
			Health: health,
			CPUPercentage: calculateCpu(stats),
			MemoryMB: memMB,
			MemoryPercentage: memPct,
		})
	}
	return  results, nil

}

func calculateCpu(s statsJSON) float64{
	cpuDelta := float64(
		s.CPUStats.CPUUsage.TotalUsage -
			s.PreCpuStats.CPUUsage.TotalUsage,
	)

	systemDelta := float64(
		s.CPUStats.SystemUsage -
			s.PreCpuStats.SystemUsage,
	)

	if cpuDelta > 0 && systemDelta > 0 {
		return (cpuDelta / systemDelta) *
			float64(s.CPUStats.OnlineCpus) * 100
	}

	return 0
}