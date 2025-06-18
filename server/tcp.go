package server

import (
	"fmt"
	"net/http"
	"sia/backend/cache"
	"sia/backend/lib"
	"sia/backend/types"

	"github.com/gorilla/mux"
)

func InitTCPServer(config *cache.Config, wsChan <-chan types.WebSocketEvent) {
	lib.Print(lib.TCP_SERVICE, "Starting TCP Server")
	r := mux.NewRouter()
	wsServer := NewWebSocketServer(config, []string{"http://localhost:3000"}) // Restrict origins in production
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsServer.ListenToWebSocket(w, r, wsChan) // Corrected order of parameters
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", lib.TCP_PORT), r)
	if err != nil {
		panic(err)
	}
}
