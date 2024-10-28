package main

import (
	"sia/backend/lib"
	"sia/backend/server"
	"sia/backend/test"
)

func main() {
	lib.InitEnvVars()
	// Run the server and client in separate goroutines
	go func() {
		server.InitUDPServer()
	}()

	go func() {
		test.RunTestUDPClient()
	}()

	// Keep the main function running to avoid exiting immediately
	select {} // Block main from exiting to let goroutines run indefinitely
}
