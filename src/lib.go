package main

import (
	"encoding/json"
	"net/http"
	"io"
	"bytes"
	"sync"
	"github.com/gorilla/websocket"
	"errors"
	"crypto/rand"
	"fmt"
)

const VERSION = "0.0.2-alpha"

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
	WriteIndex int
	ConsumerIndices map[string]int 
	Mu sync.RWMutex
}

func (mq *MessageQueue) RegisterConsumer(consumerID string) {
	mq.Mu.Lock()
	defer mq.Mu.Unlock()
	
	if _, exists := mq.ConsumerIndices[consumerID]; !exists {
		mq.ConsumerIndices[consumerID] = 0
	}
}

func (mq *MessageQueue) AddMessage(msg []byte) error {
	mq.Mu.Lock()
	defer mq.Mu.Unlock()

	minRead := mq.WriteIndex
	for _, idx := range mq.ConsumerIndices {
		if idx < minRead {
			minRead = idx
		}
	}

	if mq.WriteIndex - minRead >= QUEUE_SIZE_LIMIT {
		return errors.New("queue full - slow consumers blocking progress")
	}

	mq.Queue[mq.WriteIndex % QUEUE_SIZE_LIMIT] = msg
	mq.WriteIndex++
	return nil
}


func (mq *MessageQueue) ConsumeMessage(consumerID string) ([]byte, error) {
	mq.Mu.RLock()
	defer mq.Mu.RUnlock()

	readIdx, exists := mq.ConsumerIndices[consumerID]
	if !exists {
		return nil, errors.New(fmt.Sprintf("unregistered consumer: %s", consumerID))
	}

	available := mq.WriteIndex - readIdx
	if available <= 0 {
		return nil, errors.New("no new messages")
	}

	queueIdx := readIdx % QUEUE_SIZE_LIMIT
	msg := mq.Queue[queueIdx]

	mq.ConsumerIndices[consumerID] = readIdx + 1
	
	return msg, nil
}

func createNewMessageQueue() (*MessageQueue, error) {
	var pMessageQueue *MessageQueue = &MessageQueue{ConsumerIndices: make(map[string]int)}
	return pMessageQueue, nil
}


func openNewWsConnection(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return conn, err
	}
	return conn, err
}


func _generateUUID() (string, error) {
	uuid := make([]byte, 16)

	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0F) | 0x40
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func generateUUIDs(size int) ([]string, error) {
	uuids := make([]string, size) 
	for i := 0; i < size; i++ {
		value, err := _generateUUID()
		uuids[i] = value
		if err != nil {
			return uuids, errors.New("Error creating UUIDs.")
		}
	}
	return uuids, nil
} 