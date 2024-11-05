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
		go updateEcgChannel(ekg.Value, m, addr, ecgChan)
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
	// if output.Event != "" && output.Data != nil {
	// 	outputChan <- output
	// }
}

func updateEcgChannel(n float64, m map[string][]float64, addr string, c chan<- dsp.Signal) {
	v, f := m[addr]
	if !f {
		lib.Print(lib.UDP_SERVICE, "Not found!")
		m[addr] = make([]float64, 0)
	}

	newArr := append(v, n)
	if len(newArr) > 10000 {
		newArr = newArr[1:]
	}
	m[addr] = newArr

	c <- dsp.Signal{
		SampleRate: float64(lib.ECG_HZ),
		Signal:     newArr,
	}
}
