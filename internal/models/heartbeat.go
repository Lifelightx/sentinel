package models

type Heartbeat struct {
	
	Hostname  string `json:"hostName"`
	TimeStamp int64  `json:"timestamp"`
}

