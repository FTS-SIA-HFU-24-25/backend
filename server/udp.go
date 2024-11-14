package server

import (
	"net"
	"sia/backend/cache"
	"sia/backend/handler"
	"sia/backend/lib"
	"sia/backend/types"
)

var ecg map[string][]float64

func InitUDPServer(cache *cache.Cache, websocketChan chan<- types.WebSocketEvent, ecgChan chan types.EcgSignal) {
	ecg = make(map[string][]float64)
	conf := types.WebSocketConfigResponse{
		ChunksSize:       10,
		StartReceiveData: 0,
		FilterType:       0,
		MaxPass:          0,
		MinPass:          0,
	}
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
		handler.HandleUDPRequest(buffer, cache, conf, websocketChan, ecgChan)
	}
}
