package dockerstats

import (
	"os/exec"

	
)


func GetContainerLogs(containerID string)(string, error){
	cmd := exec.Command("docker", "logs", "--tail", "100", containerID)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func InspectContainer(containerID string)(string, error){
	cmd := exec.Command("docker", "inspect", containerID)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func RestartContainer(containerID string)(string, error){
	cmd := exec.Command("docker", "restart", containerID)
	out, err := cmd.CombinedOutput()
	return  string(out), err
}

func StopContainer(containerID string)(string, error){
	cmd := exec.Command("docker", "stop", containerID)
	out, err := cmd.CombinedOutput()
	return  string(out), err
}