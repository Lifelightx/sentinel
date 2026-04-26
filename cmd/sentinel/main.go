package main

import (
	"log"
	"os"
	"sentinel/cmd/agent"
	"sentinel/cmd/master"
)

func main(){
	role := getEnv("ROLE", "")

	if len(os.Args)> 1{
		role = os.Args[1]
	}

	switch role{
	case "master":
		log.Println("Starting master...")
		master.Run()
	case "agent":
		log.Println("Starting agent...")
		agent.Run()
	default:
		log.Fatal("Usage: sentinel [master | agent]")
	}


}


func getEnv(key, fallback string) string{
	val := os.Getenv(key)
	if val == ""{
		return  fallback
	}
	return val
}