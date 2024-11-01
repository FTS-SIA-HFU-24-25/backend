package lib

import "fmt"

const (
	UDP_SERVICE = iota
	TCP_SERVICE
	WEBSOCKET_SERVICE
)

func Print(s int, val interface{}) {
	switch s {
	case UDP_SERVICE:
		fmt.Printf("[UDP] %v \n", val)
	case TCP_SERVICE:
		fmt.Printf("[TCP] %v \n", val)
	}
}
