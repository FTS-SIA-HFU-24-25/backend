package handler

import (
	"fmt"
	"log"

	"github.com/xtaci/kcp-go/v5"
)

var (
	MAX_ERR_THRESHOLD = 5
)

func HandleSerialController(conn *kcp.UDPSession) {
	err_threshold := 0
	buf := make([]byte, 256)
	for {
		if err_threshold >= MAX_ERR_THRESHOLD {
			log.Println("Error threshold reached. Closing connection.")
			conn.Close()
			break
		}
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			err_threshold++
			continue
		}

		fmt.Println("Received ", buf[:n])

		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			err_threshold++
			continue
		}
	}
}
