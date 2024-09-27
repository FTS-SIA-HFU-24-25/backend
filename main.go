package main

import (
	"sia/backend/db"
	"sia/backend/server"
	"sia/backend/tools"
)

func init() {
	tools.LoadEnv()
}

func main() {
	tools.Log("[APP]", "Starting application...")

	// Initialize Redis and PostgreSQL
	db.InitRedis()
	db.InitDB()

	// Start the TCP server
	go func() {
		tools.Log("[TCP]", "Starting TCP server...")
		server.OpenTCPServer()
	}()

	select {}
}
