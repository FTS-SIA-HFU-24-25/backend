package main

import (
	"sia/backend/cache"
	"sia/backend/handler"
	"sia/backend/lib"
	"sia/backend/server"
	"sia/backend/types"
)

func main() {
	lib.InitEnvVars()

	websocketEventChannel := make(chan types.WebSocketEvent)
	ecgChannel := make(chan types.EcgSignal)

	cache := cache.CreateNewCache()

	go server.InitUDPServer(cache, websocketEventChannel, ecgChannel)
	go server.InitTCPServer(websocketEventChannel)
	go handler.InitECGHandler(ecgChannel, websocketEventChannel)

	select {}
}
