package main

import (
    "net/http"
    "github.com/gorilla/websocket"
	"fmt"
)

type WebSocketServer struct {
	Uuid string
	MessageQueue *MessageQueue
}

func (pWebSocketServer *WebSocketServer) Init() (error) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true 
		},
	}

	handler := func (w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
	
		for {
			msg, err := pWebSocketServer.MessageQueue.ConsumeMessage(pWebSocketServer.Uuid)
			if err != nil {
				continue
			}
			
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				continue
			}
		}
	} 

	http.HandleFunc("/"+pWebSocketServer.Uuid, handler)
	return nil
}

func createWebSocketServer(uuid string, pMessageQueue *MessageQueue) (error) {
	var webSocketServer *WebSocketServer = &WebSocketServer{
		Uuid: uuid,
		MessageQueue: pMessageQueue,
	}

	err := webSocketServer.Init()
	if err != nil {
		return err
	}

	return nil
}

func startWebSocketServers(port int) (error) {
	encodedPort := fmt.Sprintf(":%d", port)
	go http.ListenAndServe(encodedPort, nil)
	return nil
}
