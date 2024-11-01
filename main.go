package main

import (
	"sia/backend/lib"
	"sia/backend/server"
	"sia/backend/test"
	"sia/backend/types"
)

func main() {
	lib.InitEnvVars()

	websocketEventChannel := make(chan types.WebSocketEvent)
	iotEventChannel := make(chan types.IoTEvent)

	go server.InitUDPServer(iotEventChannel, websocketEventChannel)
	go test.RunTestUDPClient()
	go server.InitTCPServer(iotEventChannel, websocketEventChannel)

	select {}
}
