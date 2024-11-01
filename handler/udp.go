package handler

import (
	"fmt"
	"net"
	"sia/backend/lib"
	"sia/backend/translator"
	"sia/backend/types"
)

func HandleUDPRequest(buffer []byte, n int, addr *net.UDPAddr, outputChan chan<- types.WebSocketEvent) {
	lib.Print(lib.UDP_SERVICE, fmt.Sprintf("Received %v bytes from %v\n", n, addr))
	dataType, data := translator.TranslateUDPBinary(buffer)
	output := types.WebSocketEvent{
		Event: "",
		Data:  nil,
	}
	switch dataType {
	case types.UDP_EKG_SENSOR:
		ekg := data.(*types.EKG_SENSOR)
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("EKG sensor data: %v\n", ekg.Value))
		output.Event = "ekg-changes"
		output.Data = ekg
	case types.UDP_TEMPERATURE_SENSOR:
		temp := data.(*types.TEMPERATURE_SENSOR)
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("Temperature sensor data: %v\n", temp.Value))
		output.Event = "temp-changes"
		output.Data = temp
	case types.UDP_GPS_SERVICE:
		gps := data.(*types.GPS_SERVICE)
		latitude := gps.Latitude
		longitude := gps.Longitude
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("GPS data: %v, %v\n", latitude, longitude))
		output.Event = "gps-changes"
		output.Data = gps
	default:
		lib.Print(lib.UDP_SERVICE, "Invalid data type")
	}
	if output.Event != "" && output.Data != nil {
		outputChan <- output
	}
}
