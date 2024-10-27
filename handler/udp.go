package handler

import "net"

func HandleUDPRequest(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		panic(err)
	}
	go _handle(buffer, n, addr)
}

func _handle(buffer []byte, n int, addr *net.UDPAddr) {
	// Handle the request
}
