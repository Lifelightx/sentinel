package models

type Metrics struct{
	ServerId string `json:"serverId"`
	Hostname string `json:"hostName"`
	CPU float64 `json:"cpu"`
	RAM float64 `json:"ram"`
	Disk float64 `json:"disk"`
	Uptime uint64 `json:"uptime"`
	TimeStamp int64 `json:"timestamp"`
}