package models

const (
	ActionLogs    = "logs"
	ActionInspect = "inspect"
)

type CommandRequest struct{
	Action string `json:"action"`
	ContainerID string `json:"container_id"`
	HostName string `json:"hostName"`
	ReplyTo     string `json:"reply_to"`
}

type CommandResponse struct{
	Status string `json:"status"`
	Data any `json:"data,omitempty"`
}

type CommandResponseWrapper struct {
	HostName 	string				   `json:"hostname"`
	ContainerID string                 `json:"container_id"`
	Action      string                 `json:"action"`
	Response    CommandResponse        `json:"response"`
}

type CommandResult struct {
	HostName string
	ContainerID string
	Action      string
	Response    CommandResponse
	Timestamp   int64
}