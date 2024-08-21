package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"socket/model"
	"time"
)

func Read(client *model.Client, chatroom *model.Chat) {
	defer func() {
		client.Socket.Close()
		ChatRoom.Unregister <- client
	}()
	for {
		var msg model.SendMessage
		err := client.Socket.ReadJSON(&msg)
		msg.ID = client.ID
		if err != nil {
			log.Println("写入msg失败")
			return
		}
		log.Println(msg.Type, msg.ID, msg.ToID, string(msg.Data))
		chatroom.Broadcast <- &msg
	}
}

func Write(client *model.Client) {
	defer func() {
		client.Socket.Close()
	}()
	for {
		select {
		case msg := <-client.Send:
			msgStr, _ := json.Marshal(msg)
			client.Socket.WriteMessage(websocket.TextMessage, msgStr)
		}
	}
}

func ClientStatusCheck(client *model.Client) {
	for {
		time.Sleep(time.Second * 30)
		err := client.Socket.WriteMessage(websocket.PingMessage, nil)
		if err != nil {
			ChatRoom.Unregister <- client
			return
		}
	}
}
