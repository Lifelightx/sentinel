package models

import "time"

type ServerState struct {
	ServerId   string
	LastSeen   time.Time
	Status     string
	LastMetric Metrics
}