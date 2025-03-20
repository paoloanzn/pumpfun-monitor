package main

import (
	"encoding/json"
	"log"
	"time"
	"github.com/gorilla/websocket"
	"errors"
)

const PUMPFUN_RAYDIUM_MIGRATION_ADDRESS = "39azUYFWPz3VHgKCf3VChUwbpURdCHRxjWVowf5jUJjg"
const RAYDIUM_LIQUIDITY_POOLV4_CONTRACT_ADDRESS = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
const SOLANA_RPC_WS_URL = "ws://api.mainnet-beta.solana.com"
const TOKEN_MIGRATION_ERROR_PREFIX = "Token Migration Monitor Error: %v"

type TokenMigrationMessage struct {
	Address string `json:"address"`
	Timestamp int64 `json:"migratedOn"`
}

func startNewMigrationMonitor(pMessageQueue *MessageQueue, logger *Logger) (error) {
	logger.Info("Token Migration Monitor: Attempting connection...")
	wsConnection, err := openNewWsConnection(SOLANA_RPC_WS_URL)
	if err != nil {
		logger.Error(TOKEN_MINT_ERROR_PREFIX, err)
		return err
	}
	defer wsConnection.Close()
	logger.Success("Token Migration Monitor: Connection established.")

	wsConnection.SetPingHandler(func(data string) error {
		return wsConnection.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(10*time.Second))
	})

	// pining routine
	go func() {
        pingInterval := 15 * time.Second
        pingTimeout := 10 * time.Second
        ticker := time.NewTicker(pingInterval)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                err := wsConnection.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingTimeout))
                if err != nil {
					logger.Error(TOKEN_MIGRATION_ERROR_PREFIX, err)
                }
            }
        }
    }()

	var subMsg *JSONRpcRequest = &JSONRpcRequest{
		JsonRpc: "2.0",
		Id: 1,
		Method: "logsSubscribe",
		Params: []interface{}{
			map[string][]string{
				"mentions": []string{PUMPFUN_RAYDIUM_MIGRATION_ADDRESS},
			},
			map[string]string{
				"commitment": "finalized",
			},
		},
	}

	encodedMsg, err := json.Marshal(subMsg)
	if err != nil {
		logger.Error(TOKEN_MIGRATION_ERROR_PREFIX, err)
		return err
	}

	err = wsConnection.WriteMessage(websocket.TextMessage, encodedMsg)
	if err != nil {
		logger.Error(TOKEN_MIGRATION_ERROR_PREFIX, err)
		return err
	}

	// skip subscription confirm message
	// TODO: should ensure successful response instead
	for {
		_, _, err := wsConnection.ReadMessage()
		if err != nil {
			logger.Error(TOKEN_MIGRATION_ERROR_PREFIX, err)
			return err
		}
		break
	}

	for {
		var notificationMsg *JSONRpcNotification = &JSONRpcNotification{}
		_, message, err := wsConnection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			continue
		}

		err = json.Unmarshal(message, notificationMsg)

		if err != nil {
			logger.Warn(TOKEN_MIGRATION_ERROR_PREFIX, errors.New("Failed to parse message data, skipping."))
			continue
		}

		paramsMap, ok := notificationMsg.Params.(map[string]interface{})
		if !ok {
			logger.Warn(TOKEN_MIGRATION_ERROR_PREFIX, errors.New("Failed to parse message data, skipping."))
			continue
		}

		result, ok := paramsMap["result"].(map[string]interface{})
		if !ok {
			logger.Warn(TOKEN_MIGRATION_ERROR_PREFIX, errors.New("Failed to parse message data, skipping."))
			continue	
		}
		
		value, ok := result["value"].(map[string]interface{})
		if !ok {
			logger.Warn(TOKEN_MIGRATION_ERROR_PREFIX, errors.New("Failed to parse message data, skipping."))	
			continue	
		}
		
		signature, ok := value["signature"].(string)
		if !ok {
			logger.Warn(TOKEN_MIGRATION_ERROR_PREFIX, errors.New("Failed to parse message data, skipping."))
			continue
		}

		logs, ok := value["logs"].([]interface{})
		if !ok {
			logger.Warn(TOKEN_MIGRATION_ERROR_PREFIX, errors.New("Failed to parse message data, skipping."))
			continue
		}

		for _, txlog := range logs {
			// if this specific log is found it means its a Raydium Liquidity Pool Creation Transaction
			if txlog == "Program 675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8 invoke [1]" {
				var jsonRpcRequest *JSONRpcRequest = &JSONRpcRequest{
					JsonRpc: "2.0",
					Id: 1,
					Method: "getTransaction",
					Params: []interface{}{
						signature,
						map[string]interface{}{
							"encoding": "json", 
							"maxSupportedTransactionVersion": 0,
						},
					},
				}
				transaction, err := makeRpcHttpRequest(jsonRpcRequest)
				if err != nil {
					logger.Error(TOKEN_MIGRATION_ERROR_PREFIX, err)	
					continue
				}

				result, ok := transaction.Result.(map[string]interface{})
				if !ok {
					logger.Warn(TOKEN_MIGRATION_ERROR_PREFIX, errors.New("Failed to parse message data, skipping."))
					continue	
				}

				meta, ok := result["meta"].(map[string]interface{})
				if !ok {
					logger.Warn(TOKEN_MIGRATION_ERROR_PREFIX, errors.New("Failed to parse message data, skipping."))
					continue	
				}

				postTokenBalances, ok := meta["postTokenBalances"].([]interface{})
				if !ok {
					logger.Warn(TOKEN_MIGRATION_ERROR_PREFIX, errors.New("Failed to parse message data, skipping."))
					continue	
				}

				migratedToken, ok := postTokenBalances[1].(map[string]interface{})["mint"].(string)
				if !ok {
					logger.Warn(TOKEN_MIGRATION_ERROR_PREFIX, errors.New("Failed to parse message data, skipping."))	
					continue	
				}

				var tokenMigrationMessage *TokenMigrationMessage = &TokenMigrationMessage{
					Address: migratedToken,
					Timestamp: time.Now().UnixMilli(),
				}

				encodedMsg, err := json.Marshal(tokenMigrationMessage)
				if err != nil {
					logger.Error(TOKEN_MIGRATION_ERROR_PREFIX, err)	
					continue
				}
				pMessageQueue.AddMessage(encodedMsg)

				prettyJSON, err := json.MarshalIndent(tokenMigrationMessage, "", "  ")
				if err != nil {
					logger.Error(TOKEN_MINT_ERROR_PREFIX, "JSON formatting error:", err)
					continue
				}
		
				logger.Info("New Token Migration: %s", string(prettyJSON))
			}
		}
	}
}