package server

import (
	"net/http"
	"sia/backend/handler"
	"sia/backend/lib"
	"sia/backend/types"

	"github.com/gorilla/websocket"
)

// Allow all origins for development purposes
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust as needed for production security
	},
}

func ListenToWebSocket(w http.ResponseWriter, r *http.Request, wsChan <-chan types.WebSocketEvent) {
	lib.Print(lib.WEBSOCKET_SERVICE, "Starting websocket server")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		lib.Print(lib.WEBSOCKET_SERVICE, "Upgrade error:", err)
		return
	}
	defer c.Close()

	// Write messages to WebSocket from wsChan
	go func() {
		for data := range wsChan {
			if err := c.WriteJSON(data); err != nil {
				lib.Print(lib.WEBSOCKET_SERVICE, "WriteJSON error:", err)
				break
			}
		}
	}()

	// Read messages from WebSocket
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, "ReadMessage error:", err)
			break
		}
		// Pass the message to the handler
		handler.HandleWebsocketEvent(c, mt, message)
	}
}
