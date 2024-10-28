package main

import (
	"sia/backend/lib"
	"sia/backend/server"
	"sia/backend/test"
)

func main() {
	lib.InitEnvVars()

	go func() {
		server.InitUDPServer()
	}()
	go func() {
		test.RunTestUDPClient()
	}()

	select {}
}
