package main

import (
	"crypto/sha1"
	"log"
	"sia/backend/tools"

	"github.com/xtaci/kcp-go/v5"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	tools.Log(0, "[SERVER] ", "Starting server...")
	tools.LoadEnv()
	key := pbkdf2.Key([]byte(tools.UDP_PASSWORD), []byte(tools.UDP_SALT), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	if listener, err := kcp.ListenWithOptions(tools.UDP_ADDR, block, 10, 3); err == nil {
		// spin-up the client
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			go handleEcho(s)
		}
	} else {
		log.Fatal(err)
	}
}

// handleEcho send back everything it received
func handleEcho(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}
