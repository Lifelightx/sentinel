package models

type Metrics struct{
	
	Hostname string `json:"hostName"`
	IPv4 string `json:"ipv4"`
	CPU float64 `json:"cpu"`
	RAM float64 `json:"ram"`
	Disk float64 `json:"disk"`
	Uptime uint64 `json:"uptime"`
	TimeStamp int64 `json:"timestamp"`
}