package handler

import (
	"fmt"
	"net"
	"sia/backend/lib"
	"sia/backend/translator"
	"sia/backend/types"
)

func HandleUDPRequest(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		panic(err)
	}
	go _handle(buffer, n, addr)
}

func _handle(buffer []byte, n int, addr *net.UDPAddr) {
	lib.Print(lib.UDP_SERVICE, fmt.Sprintf("Received %v bytes from %v\n", n, addr))
	dataType, data := translator.TranslateUDPBinary(buffer)
	switch dataType {
	case types.UDP_EKG_SENSOR:
		ekg := data.(*types.EKG_SENSOR)
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("EKG sensor data: %v\n", ekg.Value))
	case types.UDP_TEMPERATURE_SENSOR:
		temp := data.(*types.TEMPERATURE_SENSOR)
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("Temperature sensor data: %v\n", temp.Value))
	case types.UDP_GPS_SERVICE:
		gps := data.(*types.GPS_SERVICE)
		latitude := gps.Latitude
		longitude := gps.Longitude
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("GPS data: %v, %v\n", latitude, longitude))
	default:
		lib.Print(lib.UDP_SERVICE, "Invalid data type")
	}
}
