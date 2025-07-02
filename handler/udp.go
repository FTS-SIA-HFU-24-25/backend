package handler

import (
	"context"
	"fmt"
	"log"
	"sia/backend/cache"
	"sia/backend/lib"
	"sia/backend/translator"
	"sia/backend/types"
)

func HandleUDPRequest(buffer []byte, cache *cache.Cache, conf *cache.Config, outputChan chan<- types.WebSocketEvent) {
	dataType, data := translator.TranslateUDPBinary(buffer)

	switch dataType {
	case types.UDP_EKG_SENSOR:
		ekg := data.(*types.EKG_SENSOR)
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("EKG sensor data: %v\n", ekg.Value))
		updateEcgChannel(ekg.Value, cache, conf, outputChan)
	case types.UDP_TEMPERATURE_SENSOR:
		if temp, ok := data.(*types.TEMPERATURE_SENSOR); ok {
			select {
			case outputChan <- types.WebSocketEvent{Event: "temp", Data: temp}:
			default:
				log.Println("[UDP_SERVICE] Dropped temp event: channel full")
			}
		}
	case types.UDP_GYRO_SENSOR:
		if gyro, ok := data.(*types.GYRO_SENSOR); ok {
			select {
			case outputChan <- types.WebSocketEvent{Event: "gyro", Data: gyro}:
			default:
				log.Println("[UDP_SERVICE] Dropped gyro event: channel full")
			}
		}
	case types.UDP_ACCEL_SENSOR:
		if accel, ok := data.(*types.ACCEL_SENSOR); ok {
			select {
			case outputChan <- types.WebSocketEvent{Event: "accel", Data: accel}:
			default:
				log.Println("[UDP_SERVICE] Dropped accel event: channel full")
			}
		}
	case types.END_CONNECTION:
		lib.Print(lib.CACHE_SERVICE, "Values cleared")
		cache.ClearValues(context.Background())
	default:
		lib.Print(lib.UDP_SERVICE, "Invalid data type")
	}
}
