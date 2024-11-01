package server

import (
	"fmt"
	"net/http"
	"sia/backend/lib"
	"sia/backend/types"

	"github.com/gorilla/mux"
)

func InitTCPServer(iotChan chan<- types.IoTEvent, wsChan <-chan types.WebSocketEvent) {
	lib.Print(lib.TCP_SERVICE, "Starting TCP Server")
	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ListenToWebSocket(w, r, iotChan, wsChan) // Corrected order of parameters
	})

	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", lib.TCP_PORT), r)
	if err != nil {
		panic(err)
	}
}
