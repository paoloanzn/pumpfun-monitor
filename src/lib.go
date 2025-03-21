package main

import (
	"encoding/json"
	"net/http"
	"io"
	"bytes"
	"sync"
	"github.com/gorilla/websocket"
	"errors"
)

type JSONRpcRequest struct {
	JsonRpc string `json:"jsonrpc"`
	Id uint64 `json:"id"`
	Method string `json:"method"`
	Params interface{} `json:"params"`
}

type JSONRpcResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result interface{} `json:"result"`
	Id uint64 `json:"id"`
}

type JSONRpcNotification struct {
	JsonRpc string `json:"jsonrpc"`
	Id string `json:"id"`
	Method string `json:"method"`
	Params interface{} `json:"params"`
	subscription string `json:"subscription"`
}

func makeRpcHttpRequest(jsonRpcRequest *JSONRpcRequest) (*JSONRpcResponse, error) {
	jsonEncoded, _ := json.Marshal(jsonRpcRequest)
	responseBody := bytes.NewBuffer(jsonEncoded)

	var jsonRpcResponse *JSONRpcResponse = &JSONRpcResponse{}

	resp, err := http.Post("https://api.mainnet-beta.solana.com", "application/json", responseBody)
	if err != nil {
		return jsonRpcResponse, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return jsonRpcResponse, err
	}

	err = json.Unmarshal(body, jsonRpcResponse)
	if err != nil {
		return jsonRpcResponse, err
	}

	return jsonRpcResponse, nil	
}

const QUEUE_SIZE_LIMIT int = 10 * 1000

type MessageQueue struct {
	Queue [QUEUE_SIZE_LIMIT][]byte
	Index int
	Mu sync.Mutex
}

func (pMessageQueue *MessageQueue) AddMessage(msg []byte) error {
	pMessageQueue.Mu.Lock()
	defer pMessageQueue.Mu.Unlock()

	if pMessageQueue.Index < QUEUE_SIZE_LIMIT {
		pMessageQueue.Queue[pMessageQueue.Index] = msg
		pMessageQueue.Index += 1

		return nil
	} else {
		return errors.New("Queue is full.")
	}
}

func (pMessageQueue *MessageQueue) ConsumeMessage() ([]byte, error) {
	pMessageQueue.Mu.Lock()
	defer pMessageQueue.Mu.Unlock()

	if pMessageQueue.Index <= 0 {
		return nil, errors.New("Queue is empty.")
	}

	pMessageQueue.Index -= 1
	msg := pMessageQueue.Queue[pMessageQueue.Index]

	return msg, nil
}

func (pMessageQueue *MessageQueue) ReadMessage() ([]byte, error){
	if pMessageQueue.Index <= 0 {
		return nil, errors.New("Queue is empty.")
	}

	return pMessageQueue.Queue[pMessageQueue.Index - 1], nil
}

func createNewMessageQueue() (*MessageQueue, error) {
	var pMessageQueue *MessageQueue = &MessageQueue{}
	return pMessageQueue, nil
}


func openNewWsConnection(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return conn, err
	}
	return conn, err
}
