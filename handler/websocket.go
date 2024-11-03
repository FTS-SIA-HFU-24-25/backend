package handler

import (
	"encoding/json"
	"sia/backend/lib"
	"sia/backend/types"
	"time"

	"github.com/gorilla/websocket"
)

type pong struct {
	Timestamp time.Time `json:"timestamp"`
}

type ping struct {
	Event int `json:"event"`
	Data  int `json:"data"`
}

func HandleWebsocketEvent(c *websocket.Conn, mt int, message []byte, outputChan chan<- types.IoTEvent) {
	var requestData ping

	err := json.Unmarshal(message, &requestData)
	if err != nil {
		lib.Print(lib.WEBSOCKET_SERVICE, err)
		return
	}

	lib.Print(lib.WEBSOCKET_SERVICE, requestData)

	if requestData.Event == types.PING {
		err := c.WriteJSON(types.WebSocketEvent{
			Event: "pong",
			Data: pong{
				Timestamp: time.Now(),
			},
		})
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			return
		}
	}
}
