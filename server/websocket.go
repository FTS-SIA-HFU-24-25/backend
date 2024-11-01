package server

import (
	"net/http"
	"sia/backend/handler"
	"sia/backend/lib"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ListenToWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		lib.Print(lib.WEBSOCKET_SERVICE, err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			break
		}

		cErr := handler.HandleWebsocketEvent(c, mt, message)
		if cErr != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, cErr.Error)
			continue
		}
	}
}
