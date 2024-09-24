package main

import (
	"crypto/sha1"
	"io"
	"log"
	"sia/backend/handler"
	"sia/backend/tools"
	"time"

	"github.com/xtaci/kcp-go/v5"
	"golang.org/x/crypto/pbkdf2"
)

func init() {
	tools.LoadEnv()
}

func main() {
	tools.Log("[SERVER]", "Starting server...")
	key := pbkdf2.Key([]byte(tools.UDP_PASSWORD), []byte(tools.UDP_SALT), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	if listener, err := kcp.ListenWithOptions(tools.UDP_ADDR, block, 10, 3); err == nil {
		// spin-up the client
		tools.Log("[SERVER]", tools.UDP_ADDR)
		go client()
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			go handler.HandleSerialController(s)
		}
	} else {
		log.Fatal(err)
	}
}

func client() {
	tools.Log("[CLIENT]", "Starting client...")
	key := pbkdf2.Key([]byte(tools.UDP_PASSWORD), []byte(tools.UDP_SALT), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	// wait for server to become ready
	time.Sleep(time.Second)

	// dial to the echo server
	if sess, err := kcp.DialWithOptions(tools.UDP_ADDR, block, 10, 3); err == nil {
		for {
			// Example: Sending a byte array
			data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
			log.Println("sent:", data)

			// Send byte array
			if _, err := sess.Write(data); err == nil {
				// Buffer to read back the same number of bytes
				buf := make([]byte, len(data))

				// Read back the data
				if _, err := io.ReadFull(sess, buf); err == nil {
					log.Println("recv:", buf)
				} else {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}

			time.Sleep(time.Second)
		}
	} else {
		log.Fatal(err)
	}
}
