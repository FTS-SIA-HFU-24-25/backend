package server

import (
	"net"
	"sia/backend/handler"
	"sia/backend/lib"
	"sia/backend/types"

	"github.com/eripe970/go-dsp-utils"
)

var ecg map[string][]float64

func InitUDPServer(iotChan <-chan types.IoTEvent, websocketChan chan<- types.WebSocketEvent, ecgChan chan dsp.Signal) {
	lib.Print(lib.UDP_SERVICE, "Starting UDP server")
	ecg = make(map[string][]float64)
	addr := net.UDPAddr{
		Port: lib.UDP_PORT,
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
		_, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}
		go handler.HandleUDPRequest(buffer, addr.String(), ecg, websocketChan, ecgChan)

		select {
		case s := <-ecgChan:
			lib.Print(lib.UDP_SERVICE, s)
		default:
			newSignal := dsp.Signal{
				SampleRate: float64(lib.ECG_HZ),
				Signal:     make([]float64, 0),
			}

			go (func() {
				ecgChan <- newSignal
			})()
		}
	}
}
