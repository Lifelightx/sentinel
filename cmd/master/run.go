package master

import (
	"encoding/json"
	"log"
	"os"
	"sentinel/internal/api"
	"sentinel/internal/broker"
	"sentinel/internal/models"
	"sentinel/internal/store"
)


func Run(){
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	mem := store.New()
	client := broker.New(natsURL)

	//sebscriber for metrics
	client.Subscribe("metrics.*", func(data []byte) {
		var m models.Metrics

		if err := json.Unmarshal(data, &m); err == nil{
			mem.SetMetrics(m.Hostname, m)
			log.Println("received:", m.Hostname, m.CPU, "%")
		}
	})
	//sebscriber for heartbeat
	client.Subscribe("heartbeat.*", func(data []byte) {
		var h models.Heartbeat
		if err := json.Unmarshal(data, &h); err == nil{
			mem.SetHeartbeat(h.Hostname, h.TimeStamp)
			log.Println("heartbeat recieved:", h.Hostname)
		}
	})
	//subscriber for containers
	client.Subscribe("containers.*", func(data []byte) {
	var payload models.ContainerPayload
	if err := json.Unmarshal(data, &payload); err == nil{
		mem.SetContainers(payload.HostName, payload.Containers)
	}
	log.Println("container stats received:", payload.HostName)
})

	client.Subscribe("reply.*", func(data []byte) {
		var resp models.CommandResponseWrapper
		
		if err := json.Unmarshal(data, &resp); err == nil{
			mem.SetCommandResult(resp.HostName, resp.ContainerID,resp.Action, resp.Response)
			log.Println("command response stored for: ", resp.HostName, resp.ContainerID, resp.Action)
		}

	})
	
	log.Println("master started at :8080")
	log.Fatal(api.Start(":8080",mem, client))
}

func getEnv(key , fallback string) string{
	val := os.Getenv(key)
	if val == ""{
		return fallback
	}
	return  val
}