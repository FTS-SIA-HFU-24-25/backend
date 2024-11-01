package server

import (
	"net/http"
	"sia/backend/handler"
	"sia/backend/lib"
	"sia/backend/types"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ListenToWebSocket(w http.ResponseWriter, r *http.Request, iotChan chan<- types.IoTEvent, wsChan <-chan types.WebSocketEvent) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		lib.Print(lib.WEBSOCKET_SERVICE, err)
		return
	}
	defer c.Close()

	go (func() {
		for data := range wsChan {
			err := c.WriteJSON(data)
			if err != nil {
				break
			}
		}
	})()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			lib.Print(lib.WEBSOCKET_SERVICE, err)
			break
		}

		go handler.HandleWebsocketEvent(c, mt, message, iotChan)
	}
}
