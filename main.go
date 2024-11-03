package main

import (
	"sia/backend/handler"
	"sia/backend/lib"
	"sia/backend/server"
	"sia/backend/types"

	"github.com/eripe970/go-dsp-utils"
)

const (
	EcgHZ float64 = 100
)

func main() {
	lib.InitEnvVars()

	websocketEventChannel := make(chan types.WebSocketEvent)
	iotEventChannel := make(chan types.IoTEvent)
	ecgChannel := make(chan dsp.Signal)

	go server.InitUDPServer(iotEventChannel, websocketEventChannel, ecgChannel)
	go server.InitTCPServer(iotEventChannel, websocketEventChannel)
	go handler.InitECGHandler(ecgChannel, websocketEventChannel)

	select {}
}
