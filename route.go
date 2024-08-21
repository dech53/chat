package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"socket/model"
	"socket/utils"
	"strconv"
	"sync"
	"time"
)

var ctx = context.Background()
var ChatRoom = model.Chat{
	Clients:    make(map[string]*model.Client),
	Broadcast:  make(chan *model.SendMessage),
	Register:   make(chan *model.Client),
	Unregister: make(chan *model.Client),
	Reply:      make(chan *model.Client),
}

// 用户唯一ID
var id = 0
var mu sync.Mutex

func Wsfc(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	go ClientTask()

	newClient := &model.Client{
		ID:     strconv.Itoa(id),
		Socket: conn,
		Send:   make(chan *model.SendMessage),
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{"localhost:9092"},
			Topic:          strconv.Itoa(id),
			CommitInterval: 1 * time.Second,
			GroupID:        strconv.Itoa(id), //定义消费组名字
			StartOffset:    kafka.FirstOffset,
		}),
	}
	go ClientStatusCheck(newClient)
	go Read(newClient, &ChatRoom)
	go Write(newClient)
	ChatRoom.Register <- newClient
	id++
	//阻塞主线程
	select {}
}
func ClientTask() {
	for {
		select {
		case client := <-ChatRoom.Register:
			mu.Lock()
			log.Printf("用户%v已建立新的连接\n", client.ID)
			ChatRoom.Clients[client.ID] = client
			mu.Unlock()
			replyMsg := &model.ReplyMsg{
				From:    "Server",
				Code:    1,
				Content: "已连接至服务器",
			}
			msg, _ := json.Marshal(replyMsg)
			client.Socket.WriteMessage(websocket.TextMessage, msg)
			go utils.ReadKafka(&ctx, client)
		case client := <-ChatRoom.Unregister:
			mu.Lock()
			log.Printf("用户%v断开连接\n", client.ID)
			if _, ok := ChatRoom.Clients[client.ID]; ok {
				replyMsg := &model.ReplyMsg{
					From:    "Server",
					Code:    2,
					Content: "连接已断开",
				}
				msg, _ := json.Marshal(replyMsg)
				client.Socket.WriteMessage(websocket.TextMessage, msg)
				close(client.Send)
				delete(ChatRoom.Clients, client.ID)
				client.Reader.Close()
			}
			mu.Unlock()
		case message := <-ChatRoom.Broadcast:
			mu.Lock()
			if message.Type < -1 || (message.ID == message.ToID) {
				log.Println("消息类型不合规")
			} else if message.Type == -1 {
				for _, client := range ChatRoom.Clients {
					utils.WriteKafka(message, ctx, client.ID)
				}
			} else {
				client1 := ChatRoom.Clients[message.ToID]
				client2 := ChatRoom.Clients[message.ID]
				if client1 == nil {
					utils.WriteKafka(message, ctx, client2.ID)
					utils.WriteKafka(message, ctx, message.ToID)
				} else {
					utils.WriteKafka(message, ctx, client2.ID)
					utils.WriteKafka(message, ctx, client1.ID)
				}

			}
			mu.Unlock()
		}
	}
}
