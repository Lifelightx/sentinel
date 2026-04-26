package store

import (
	"sentinel/internal/models"
	"sync"
	"time"
)

type MemoryStore struct{
	mu sync.RWMutex
	metrics map[string]models.Metrics
	lastSeen map[string]int64
	containers map[string][]models.ContainerInfo
}

func New() *MemoryStore{
	return &MemoryStore{
		metrics: make(map[string]models.Metrics),
		lastSeen: make(map[string]int64),
		containers: make(map[string][]models.ContainerInfo),
	}
}

func (s *MemoryStore) SetMetrics(id string, metric models.Metrics){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics[id] = metric
	s.lastSeen[id] = metric.TimeStamp
}

func (s *MemoryStore) SetHeartbeat(id string, ts int64){
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastSeen[id] = ts
}
func (s *MemoryStore) SetContainers(id string, data[]models.ContainerInfo){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.containers[id] = data
}

func (s *MemoryStore) GetAll() []map[string]any{
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []map[string]any{}
	now := time.Now().Unix()

	for id, metric := range s.metrics{
		status :="online"

		if now - s.lastSeen[id] > 30{
			status = "offline"
		}
		out = append(out, map[string]any{
			"serverId":id,
			"hostname":metric.Hostname,
			"cpu":metric.CPU,
			"ram":metric.RAM,
			"disk":metric.Disk,
			"lastseen":s.lastSeen[id],
			"status":status,
		})
	}
	return out
}

func (s *MemoryStore) GetContainers(id string) []models.ContainerInfo{
	s.mu.RLock()
	defer s.mu.RUnlock()
	return  s.containers[id]
}