package main

import (
	"sia/backend/cache"
	"sia/backend/lib"
	"sia/backend/server"
	"sia/backend/types"
)

func main() {
	lib.InitEnvVars()

	websocketEventChannel := make(chan types.WebSocketEvent, 1024)

	cache, config, err := cache.CreateNewCache()
	if err != nil {
		panic(err)
	}

	go server.InitUDPServer(cache, config, websocketEventChannel)
	go server.InitTCPServer(config, websocketEventChannel)

	select {}
}
