package main

import (
	"github.com/gorilla/websocket"
	"encoding/json"
)

type NewTokenMessage struct {
	Signature       string  `json:"signature"`
	Mint            string  `json:"mint"`
	TraderPublicKey string  `json:"traderPublicKey"`
	TxType          string  `json:"txType"`
	InitialBuy      float64 `json:"initialBuy"`
	SolAmount       float64 `json:"solAmount"`
	BondingCurveKey string  `json:"bondingCurveKey"`
	VTokensInBondingCurve float64 `json:"vTokensInBondingCurve"`
	VSolInBondingCurve float64 `json:"vSolInBondingCurve"`
	MarketCapSol     float64 `json:"marketCapSol"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Uri              string  `json:"uri"`
	Pool             string  `json:"pool"`
}

const PUMPFUN_PORTAL_WS_API_URL = "wss://pumpportal.fun/api/data"
const TOKEN_MINT_ERROR_PREFIX = "Token Mint Monitor Error: %v"

func startNewTokenMintMonitor(pMessageQueue *MessageQueue, logger *Logger) (error) {
	logger.Info("Token Minting Monitor: Attempting connection...")
	wsConnection, err := openNewWsConnection(PUMPFUN_PORTAL_WS_API_URL)
	if err != nil {
		logger.Error(TOKEN_MINT_ERROR_PREFIX, err)
		return err
	}
	defer wsConnection.Close()
	logger.Success("Token Minting Monitor: Connection established.")

	subMsg := `{"method": "subscribeNewToken"}`

	err = wsConnection.WriteMessage(websocket.TextMessage, []byte(subMsg))
	if err != nil {
		logger.Error(TOKEN_MINT_ERROR_PREFIX, err)
		return err
	}

	// skip subscription confirm message
	for {
		_, _, err := wsConnection.ReadMessage()
		if err != nil {
			logger.Error(TOKEN_MINT_ERROR_PREFIX, err)
			return err
		}
		break
	}

	for {
		_, message, err := wsConnection.ReadMessage()
		if err != nil {
			logger.Error(TOKEN_MINT_ERROR_PREFIX, err)
			continue
		}
		var newTokenMessage *NewTokenMessage = &NewTokenMessage{}
		err = json.Unmarshal(message, newTokenMessage)
		if err != nil {
			logger.Error(TOKEN_MINT_ERROR_PREFIX, err)
			continue
		}

		prettyJSON, err := json.MarshalIndent(newTokenMessage, "", "  ")
		if err != nil {
			logger.Error(TOKEN_MINT_ERROR_PREFIX, "JSON formatting error:", err)
			continue
		}

		logger.Info("New token: %s", string(prettyJSON))
		pMessageQueue.AddMessage(message)
	}
}


