package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"sentinel/internal/broker"
	"sentinel/internal/dockerstats"
	"sentinel/internal/metrics"
	"sentinel/internal/models"
)


func Run(){
	natsUrl := getEnv("NATS_URL", "nats://localhost:4222")
	hostName := getHostName()

	client := broker.New(natsUrl)
	log.Println("Agent Started:", hostName)
	client.Subscribe("commands."+hostName, func(data []byte) {
		log.Println("Agent listend for that....")
		var cmd models.CommandRequest
		err := json.Unmarshal(data, &cmd)
		if err != nil{
			log.Println("Invalid command", err)
			return
		}
		var result any
		var execErr error
		switch cmd.Action{
		case models.ActionLogs:
			result, execErr = dockerstats.GetContainerLogs(cmd.ContainerID)
		case models.ActionInspect:
			result, execErr = dockerstats.InspectContainer(cmd.ContainerID)
		case models.Restart:
			result, execErr = dockerstats.RestartContainer(cmd.ContainerID)
		case models.Stop:
			result, execErr = dockerstats.StopContainer(cmd.ContainerID)

		default:
			execErr = fmt.Errorf("unknown action: %s", cmd.Action)
		}
		//build response
		resp := models.CommandResponseWrapper{
			HostName: hostName,
			ContainerID: cmd.ContainerID,
			Action: cmd.Action,
			Response: models.CommandResponse{
				Status: "success",
				Data: result,
			},

		}
		if execErr != nil{
			resp.Response.Status = "error"
			resp.Response.Data = execErr.Error()
		}
		//send response back
		
		err = client.Publish(cmd.ReplyTo, resp)
		if err != nil{
			fmt.Println("failed send response: ", err)
		}

	})
	
	go startHeartbeatLoop(client, hostName)
	go dockerStatsContainerLoop(client, hostName)
	go startMetricsLoop(client, hostName)
	select {}
	
	
}

func startMetricsLoop(client *broker.Client, serverId string){
	for{
		data, err := metrics.Collect()
		if err == nil{
			publish(client, "metrics."+serverId, data,  "metrics sent: ")
			// log.Println(serverId)
		}
		time.Sleep(5 * time.Second)
	}
}

func startHeartbeatLoop(client *broker.Client, serverId string){
	for{
		data := models.Heartbeat{
			Hostname: serverId,
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

func dockerStatsContainerLoop(client *broker.Client, hostName string){
	for{
		data, err := dockerstats.Collect()
	if err == nil{
		payload := models.ContainerPayload{
			HostName: hostName,
			TimeStamp: time.Now().Unix(),
			Containers: data,
		}
		publish(client, "containers."+hostName, payload, "containres data sent")
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


func getHostName() string{
	hostname, err := os.Hostname()
	if err == nil && hostname != "" {
		return hostname
	}
	return "Unknown-server"
}
