package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func main() {
	dialer := &websocket.Dialer{}
	header := http.Header{"Cookie": []string{"name=dech53"}}
	conn, resp, err := dialer.Dial("ws://localhost:8080/ws", header)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("拨号错误")
	}
	for key, values := range resp.Header {
		fmt.Printf("%s:%s\n", key, values[0])
	}
	defer conn.Close()
}
