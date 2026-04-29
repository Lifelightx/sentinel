package store

import (
	"sentinel/internal/models"
	"sort"
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

func (s *MemoryStore) GetAll() []models.ServerView{
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []models.ServerView{}
	now := time.Now().Unix()

	for id, metric := range s.metrics{
		status :="online"
		score, count := s.calCulateAlertScore(metric, id)
		if now - s.lastSeen[id] > 30{
			status = "offline"
		}
		out = append(out, models.ServerView{
			ServerID:id,
			Hostname:metric.Hostname,
			CPU:metric.CPU,
			RAM:metric.RAM,
			Disk:metric.Disk,
			IPv4: metric.IPv4,
			LastSeen:s.lastSeen[id],
			Status:status,
			AlertCount: count,
			AlertScore: score,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].AlertScore == out[j].AlertScore {
			return out[i].AlertCount > out[j].AlertCount
		}
		return out[i].AlertScore > out[j].AlertScore
	})

	return out
}

func (s *MemoryStore) GetByID(id string) (models.ServerView, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	metric, ok := s.metrics[id]
	if !ok {
		return models.ServerView{}, false
	}

	status := "online"
	if time.Now().Unix()-s.lastSeen[id] > 30 {
		status = "offline"
	}
	
	
	return models.ServerView{
		ServerID: id,
		Hostname: metric.Hostname,
		CPU:      metric.CPU,
		RAM:      metric.RAM,
		Disk:     metric.Disk,
		IPv4: 	  metric.IPv4,
		LastSeen: s.lastSeen[id],
		Status:   status,
	}, true
}

func (s *MemoryStore) GetContainers(id string) []models.ContainerInfo{
	s.mu.RLock()
	defer s.mu.RUnlock()
	return  s.containers[id]
}

func (s *MemoryStore) calCulateAlertScore(metric models.Metrics, id string)(int, int){
	score :=0
	count :=0
	if metric.CPU > 95{
		score += 15
		count++
	}
	if metric.RAM > 90{
		score += 20
		count++
	}
	if metric.Disk > 90{
		score += 30
		count++
	}
	for _, c := range s.containers[id]{
		if c.Health == "unhealthy"{
			score += 40
			count++
		}
		if c.State == "exited"{
			score += 10
			count ++
		}

	}
	return score, count
}