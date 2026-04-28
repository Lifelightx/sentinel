package models

import "time"

type ServerState struct {
	ServerId   string
	LastSeen   time.Time
	Status     string
	LastMetric Metrics
}

type ServerView struct {
    ServerID    string
    Hostname    string
    CPU         float64
    RAM         float64
    Disk        float64
    Status      string
    LastSeen    int64
    AlertScore  int
    AlertCount  int
}