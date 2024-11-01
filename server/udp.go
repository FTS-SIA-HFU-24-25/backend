package server

import (
	"net"
	"sia/backend/handler"
	"sia/backend/lib"
	"sia/backend/types"
)

func InitUDPServer(iotChan <-chan types.IoTEvent, websocketChan chan<- types.WebSocketEvent) {
	lib.Print(lib.UDP_SERVICE, "Starting UDP server")
	addr := net.UDPAddr{
		Port: int(lib.UDP_PORT),
		IP:   net.ParseIP(lib.UDP_ADDR),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go (func() {
		for data := range iotChan {
			response := []byte{byte(data.Type)}
			response = append(response, data.Data...)
			conn.WriteToUDP(response, &addr)
		}
	})()

	for {
		buffer := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}
		go handler.HandleUDPRequest(buffer, n, addr, websocketChan)
	}
}
