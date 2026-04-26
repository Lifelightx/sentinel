package main

import (
	"log"
	"os"
	"time"

	"sentinel/internal/broker"
	"sentinel/internal/dockerstats"
	"sentinel/internal/metrics"
	"sentinel/internal/models"
)


func main(){
	natsUrl := getEnv("NATS_URL", "nats://localhost:4222")
	serverId := getEnv("SERVER_ID", "server-1")

	client := broker.New(natsUrl)
	log.Println("Agent Started:", serverId)
	
	go startHeartbeatLoop(client, serverId)
	go dockerStatsContainerLoop(client, serverId)
	startMetricsLoop(client, serverId)
	
}

func startMetricsLoop(client *broker.Client, serverId string){
	for{
		data, err := metrics.Collect(serverId)
		if err == nil{
			publish(client, "metrics."+serverId, data,  "metrics sent")
		}
		time.Sleep(5 * time.Second)
	}
}

func startHeartbeatLoop(client *broker.Client, serverId string){
	for{
		data := models.Heartbeat{
			ServerId: serverId,
			TimeStamp: time.Now().Unix(),
		}
		publish(client, "heartbeat."+serverId, data,"heartbeat sent")
		time.Sleep(10* time.Second)
	}
}



func publish(client *broker.Client, subject string, data any, message string){
	err := client.Publish(subject, data)
	if err != nil{
		log.Println("publish error: ", err)
	}
	log.Println("published: ",message)
}

func dockerStatsContainerLoop(client *broker.Client, serverId string){
	for{
		data, err := dockerstats.Collect()
	if err == nil{
		payload := models.ContainerPayload{
			ServerId: serverId,
			TimeStamp: time.Now().Unix(),
			Containers: data,
		}
		publish(client, "containers."+serverId, payload, "containres data sent")
	}else{
		log.Println("container stats error", err)
	}

	time.Sleep(15 * time.Second)
	}
}

func getEnv(key, fallback string) string{
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}