package server

import (
	"encoding/json"
	"net"
	"sia/backend/cache"
	"sia/backend/handler"
	"sia/backend/lib"
	"sia/backend/types"
)

var ecg map[string][]float64

func InitUDPServer(cache *cache.Cache, iotChan <-chan types.IoTEvent, websocketChan chan<- types.WebSocketEvent, ecgChan chan types.EcgSignal) {
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

	go (func() {
		for data := range iotChan {
			if data.Type == 255 {
				var wsConfMsg types.WebSocketConfigResponse
				err := json.Unmarshal(data.Data, &wsConfMsg)
				if err != nil {
					continue
				}
				conf = wsConfMsg
				continue
			}
			response := []byte{byte(data.Type)}
			response = append(response, data.Data...)
			conn.WriteToUDP(response, &addr)
		}
	})()

	for {
		buffer := make([]byte, 1024)
		_, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			lib.Print(lib.UDP_SERVICE, err)
			continue
		}
		go handler.HandleUDPRequest(buffer, cache, conf, websocketChan, ecgChan)
	}
}
