package server

import (
	"net"
	"sia/backend/cache"
	"sia/backend/handler"
	"sia/backend/lib"
	"sia/backend/types"
)

var (
	ecg map[string][]float64
)

func InitUDPServer(cache *cache.Cache, config *cache.Config, websocketChan chan<- types.WebSocketEvent) {
	ecg = make(map[string][]float64)
	addr := net.UDPAddr{
		Port: lib.UDP_PORT,
		IP:   net.ParseIP(lib.UDP_ADDR),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	lib.Print(lib.UDP_SERVICE, "Starting UDP server: ", addr.String())

	for {
		buffer := make([]byte, 1024)
		_, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			lib.Print(lib.UDP_SERVICE, err)
			continue
		}
		handler.HandleUDPRequest(buffer, cache, config, websocketChan)
	}
}
