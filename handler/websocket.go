package handler

import (
	"encoding/json"
	"sia/backend/lib"
	"sia/backend/types"

	"github.com/gorilla/websocket"
)

func HandleWebsocketEvent(c *websocket.Conn, mt int, message []byte) *lib.CustomError {
	var requestData types.WebSocketRequest

	err := json.Unmarshal(message, &requestData)
	if err != nil {
		return lib.CreateError(lib.WEBSOCKET_SERVICE, err)
	}

	switch requestData.Event {
	case 0: // message

	}

	return nil
}
