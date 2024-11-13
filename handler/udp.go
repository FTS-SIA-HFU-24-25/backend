package handler

import (
	"context"
	"fmt"
	"sia/backend/cache"
	"sia/backend/lib"
	"sia/backend/translator"
	"sia/backend/types"

	"github.com/eripe970/go-dsp-utils"
)

func HandleUDPRequest(buffer []byte, cache *cache.Cache, conf types.WebSocketConfigResponse, outputChan chan<- types.WebSocketEvent, ecgChan chan<- types.EcgSignal) {
	dataType, data := translator.TranslateUDPBinary(buffer)
	output := types.WebSocketEvent{
		Event: "",
		Data:  nil,
	}
	switch dataType {
	case types.UDP_EKG_SENSOR:
		ekg := data.(*types.EKG_SENSOR)
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("EKG sensor data: %v\n", ekg.Value))
		updateEcgChannel(ekg.Value, cache, conf, ecgChan)
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
	case types.END_CONNECTION:
		lib.Print(lib.CACHE_SERVICE, "Values cleared")
		cache.ClearValues(context.Background())
	default:
		lib.Print(lib.UDP_SERVICE, "Invalid data type")
	}
	// if output.Event != "" && output.Data != nil {
	// 	outputChan <- output
	// }
}

func updateEcgChannel(n float64, cache *cache.Cache, conf types.WebSocketConfigResponse, c chan<- types.EcgSignal) {
	arr, err := cache.AddIndexToEcg(context.TODO(), n)
	if err != nil {
		lib.Print(lib.UDP_SERVICE, err)
		return
	}
	length := len(*arr)

	lib.Print(lib.UDP_SERVICE, arr, length)

	if length%conf.ChunksSize > 0 {
		return
	}

	c <- types.EcgSignal{
		Signal: dsp.Signal{
			SampleRate: float64(lib.ECG_HZ),
			Signal:     *arr,
		},
		ChunksSize: conf.ChunksSize,
	}
}
