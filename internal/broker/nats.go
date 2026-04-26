package broker

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)



type Client struct{
	Conn *nats.Conn
}

func New(url string) *Client{
	nc, err := nats.Connect(url)
	if err != nil{
		log.Fatal(err)
	}
	return &Client{Conn: nc}
}

func (c *Client) Publish(subject string, data any) error{
	payload, err := json.Marshal(data)
	if err != nil{
		return err
	}
	return c.Conn.Publish(subject, payload)
}

func (c *Client) Subscribe(subject string, handler func([]byte)) error{
	_, err := c.Conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	return err
}