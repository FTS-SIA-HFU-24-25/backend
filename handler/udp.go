package handler

import (
	"fmt"
	"sia/backend/lib"
	"sia/backend/translator"
	"sia/backend/types"

	"github.com/eripe970/go-dsp-utils"
)

func HandleUDPRequest(buffer []byte, addr string, m map[string][]float64, outputChan chan<- types.WebSocketEvent, ecgChan chan<- dsp.Signal) {
	dataType, data := translator.TranslateUDPBinary(buffer)
	output := types.WebSocketEvent{
		Event: "",
		Data:  nil,
	}
	switch dataType {
	case types.UDP_EKG_SENSOR:
		ekg := data.(*types.EKG_SENSOR)
		// lib.Print(lib.UDP_SERVICE, fmt.Sprintf("EKG sensor data: %v\n", ekg.Value))
		output.Event = "ekg-changes"
		output.Data = ekg
		v, f := m[addr]
		if !f {
			lib.Print(lib.UDP_SERVICE, "Not found!")
			m[addr] = make([]float64, 0)
		}
		newArr := append(v, float64(ekg.Value))
		m[addr] = newArr
		go updateEcgChannel(float64(ekg.Value), newArr, ecgChan)
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

func updateEcgChannel(v float64, arr []float64, c chan<- dsp.Signal) {
	c <- dsp.Signal{
		SampleRate: 100,
		Signal:     arr,
	}
}
