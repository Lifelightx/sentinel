package models

type Heartbeat struct {
	ServerId  string `json:"serverId"`
	Hostname  string `json:"hostName"`
	TimeStamp int64  `json:"timestamp"`
}

