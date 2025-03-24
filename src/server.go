package main

import (
    "net/http"
    "github.com/gorilla/websocket"
	"fmt"
	"sync"
)

type WebSocketServer struct {
	Uuid string
	MessageQueue *MessageQueue
	Clients map[*websocket.Conn]struct{}
	ClientsBuffers map[*websocket.Conn]chan []byte
	ClientsMutex sync.Mutex
}

func (pWebSocketServer *WebSocketServer) RunBroadcaster() (error) {
	consumerCh, err := pWebSocketServer.MessageQueue.GetConsumerChannel(pWebSocketServer.Uuid)
	if err != nil {
		return err
	}
	
	for msg := range consumerCh {
		pWebSocketServer.ClientsMutex.Lock() 
		for _, buffer := range pWebSocketServer.ClientsBuffers {
			select {
			case buffer <- msg:
			default:
			}
		}
		pWebSocketServer.ClientsMutex.Unlock()
	}

	return nil
}

func (pWebSocketServer *WebSocketServer) Init() (error) {
	var upgrader = websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }

	pWebSocketServer.Clients = make(map[*websocket.Conn]struct{})
	pWebSocketServer.ClientsBuffers = make(map[*websocket.Conn]chan []byte)

	go pWebSocketServer.RunBroadcaster()

    handler := func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            return
        }
        defer conn.Close()

		pWebSocketServer.ClientsMutex.Lock()
		pWebSocketServer.Clients[conn] = struct{}{}
		pWebSocketServer.ClientsBuffers[conn] = make(chan []byte, 100)
		pWebSocketServer.ClientsMutex.Unlock()

		// unregister client on disconnection
		defer func() {
			pWebSocketServer.ClientsMutex.Lock()
			delete(pWebSocketServer.Clients, conn)
			if buffer, exists := pWebSocketServer.ClientsBuffers[conn]; exists {
				close(buffer)
				delete(pWebSocketServer.ClientsBuffers, conn)
			}
			pWebSocketServer.ClientsMutex.Unlock()	
			conn.Close()
		}()
        
        for msg := range pWebSocketServer.ClientsBuffers[conn] {
            if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
                break  // Client disconnected
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
