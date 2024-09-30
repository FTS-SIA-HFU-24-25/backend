package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sia/backend/handlers"
	"sia/backend/tools"

	"github.com/xtaci/kcp-go/v5"
)

func InitUDPServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	block, _ := kcp.NewNoneBlockCrypt([]byte{})

	listener, err := kcp.ListenWithOptions(tools.UDP_ADDR, block, 10, 3)
	if err != nil {
		tools.Log("[UDP]", "Failed to listen to port")
		panic(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				tools.Log("[UDP]", "Server closed")
				listener.Close()
				return
			default:
				conn, err := listener.AcceptKCP()
				if err != nil {
					tools.Log("[UDP]", "Failed to accept connection")
					continue
				}
				go handlers.HandleUDPConn(conn)
			}
		}
	}()

	<-sigs
	cancel()
	time.Sleep(1 * time.Second)
	tools.Log("[UDP]", "Server shut down")
}
