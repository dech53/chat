package model

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan *SendMessage
	Reader *kafka.Reader
}
type SendMessage struct {
	Type int             `json:"type"`
	ID   string          `json:"id"`
	ToID string          `json:"to_id"`
	Data json.RawMessage `json:"data"`
}
type Chat struct {
	Clients    map[string]*Client
	Broadcast  chan *SendMessage
	Register   chan *Client
	Unregister chan *Client
	Reply      chan *Client
}
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}
