package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"sia/backend/cache"
	"sia/backend/lib"
	"sia/backend/types"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketIncomingEvent struct {
	Event int    `json:"event"`
	Data  string `json:"data"`
}

func HandleWebsocketEvent(c *websocket.Conn, mt int, message []byte, config *cache.Config) {
	var requestData WebSocketIncomingEvent

	err := json.Unmarshal(message, &requestData)
	if err != nil {
		lib.Print(lib.WEBSOCKET_SERVICE, err)
		return
	}

	lib.Print(lib.WEBSOCKET_SERVICE, fmt.Sprintf("Ping > %d", time.Now().Unix()))

	switch requestData.Event {
	case types.PING:
		conf, err := config.GetConfig(context.Background())
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			return
		}
		err = c.WriteJSON(types.WebSocketEvent{
			Event: "pong",
			Data:  conf,
		})
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			return
		}
	case types.CONFIG_UPDATE:
		var conf types.WebSocketConfigResponse
		err := json.Unmarshal([]byte(requestData.Data), &conf)
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			return
		}
	case types.START_ECG:
		conf, err := config.GetConfig(context.Background())
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			return
		}
		conf.StartReceiveData = 1
		err = config.Set(context.Background(), "config", *conf)
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			return
		}
	case types.STOP_ECG:
		conf, err := config.GetConfig(context.Background())
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			return
		}
		conf.StartReceiveData = 0
		err = config.Set(context.Background(), "config", *conf)
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			return
		}
	case types.SPECTRUM_UPDATE:
		conf, err := config.GetConfig(context.Background())
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			return
		}
		conf.SpectrumUpdateRequest = 1
		err = config.Set(context.Background(), "config", *conf)
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			return
		}
	}
}
