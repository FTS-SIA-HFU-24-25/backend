package server

import (
	"net"
	"sia/backend/handler"
	"sia/backend/lib"
)

func InitUDPServer() {
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

	for {
		handler.HandleUDPRequest(conn)
	}
}
