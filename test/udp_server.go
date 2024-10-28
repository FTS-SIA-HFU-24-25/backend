package test

import (
	"log"
	"math/rand"
	"net"
	"sia/backend/lib"
	"sia/backend/types"
	"time"
)

func sendData(conn *net.UDPConn, dataType int, data []byte) {
	packet := append([]byte{byte(dataType)}, data...)
	_, err := conn.Write(packet)
	if err != nil {
		log.Printf("Failed to send data: %v", err)
	}
}

func RunTestUDPClient() {
	lib.Print(lib.UDP_SERVICE, "Starting UDP client")
	serverAddr := net.UDPAddr{
		Port: int(lib.UDP_PORT),
		IP:   net.ParseIP(lib.UDP_ADDR),
	}
	conn, err := net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		lib.Print(lib.UDP_SERVICE, "Failed to connect to server")
	}
	defer conn.Close()

	for {
		dataType := rand.Intn(3)
		switch dataType {
		case types.UDP_EKG_SENSOR:
			data := byte(rand.Intn(256))
			sendData(conn, dataType, []byte{data})
			lib.Print(lib.UDP_SERVICE, "Sent EKG sensor data")

		case types.UDP_TEMPERATURE_SENSOR:
			temp := rand.Float64()*100 - 50 // Temperature range -50 to +50
			sendData(conn, dataType, lib.Float64ToBytes(temp))
			lib.Print(lib.UDP_SERVICE, "Sent temperature sensor data")

		case types.UDP_GPS_SERVICE:
			latitude := rand.Float64()*180 - 90   // Latitude range -90 to +90
			longitude := rand.Float64()*360 - 180 // Longitude range -180 to +180
			data := append(lib.Float64ToBytes(latitude), lib.Float64ToBytes(longitude)...)
			sendData(conn, dataType, data)
			lib.Print(lib.UDP_SERVICE, "Sent GPS data")
		default:
			lib.Print(lib.UDP_SERVICE, "Invalid data type")
		}
		time.Sleep(1 * time.Second)
	}
}
