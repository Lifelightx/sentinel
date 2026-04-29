package models

type CommandRequest struct{
	Action string `json:"action"`
	ContainerID string `json:"container_id"`
	ServerID string `json:"server_id"`
	ReplyTo     string `json:"reply_to"`
}

type CommandResponse struct{
	Status string `json:"status"`
	Data any `json:"data,omitempty"`
}
