package handler

import (
	"encoding/json"
	"sia/backend/lib"
	"sia/backend/types"

	"github.com/gorilla/websocket"
)

func HandleWebsocketEvent(c *websocket.Conn, mt int, message []byte, outputChan chan<- types.IoTEvent) {
	var requestData types.WebSocketRequest

	err := json.Unmarshal(message, &requestData)
	if err != nil {
		lib.Print(lib.WEBSOCKET_SERVICE, err)
		return
	}

	switch requestData.Event {
	case 0: // message

	}
}
