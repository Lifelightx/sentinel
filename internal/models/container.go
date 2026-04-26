package models


type ContainerInfo struct{
	ID string `json:"id"`
	Name string `json:"name"`
	State string `json:"state"`
	Status string `json:"status"`
	Health string `json:"health"`
	CPUPercentage float64 `json:"cpuPercentage"`
	MemoryMB float64 `json:"memoryMB"`
	MemoryPercentage float64 `json:"memoryPercentage"`
}

type ContainerPayload struct{
	ServerId string `json:"serverId"`
	TimeStamp int64 `json:"timestamp"`
	Containers []ContainerInfo `json:"containers"`
}