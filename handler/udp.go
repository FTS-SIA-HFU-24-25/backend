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

func HandleUDPRequest(buffer []byte, cache *cache.Cache, conf *cache.Config, outputChan chan<- types.WebSocketEvent) {
	dataType, data := translator.TranslateUDPBinary(buffer)

	switch dataType {
	case types.UDP_EKG_SENSOR:
		ekg := data.(*types.EKG_SENSOR)
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("EKG sensor data: %v\n", ekg.Value))
		updateEcgChannel(ekg.Value, cache, conf, outputChan)
	case types.UDP_TEMPERATURE_SENSOR:
		temp := data.(*types.TEMPERATURE_SENSOR)
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("Temperature sensor data: %v\n", temp.Value))
		outputChan <- types.WebSocketEvent{
			Event: "temp",
			Data: temp,
		}
	case types.UDP_GYRO_SENSOR:
		gyro := data.(*types.GYRO_SENSOR)
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("Gyro sensor data: %v\n", gyro))
		outputChan <- types.WebSocketEvent{
			Event: "gyro",
			Data: gyro,
		}
	case types.UDP_ACCEL_SENSOR:
		accel := data.(*types.ACCEL_SENSOR)
		lib.Print(lib.UDP_SERVICE, fmt.Sprintf("Accel sensor data: %v\n", accel))
		outputChan <- types.WebSocketEvent{
			Event: "accel",
			Data: accel,
		}
	case types.END_CONNECTION:
		lib.Print(lib.CACHE_SERVICE, "Values cleared")
		cache.ClearValues(context.Background())
	default:
		lib.Print(lib.UDP_SERVICE, "Invalid data type")
	}
}

func updateEcgChannel(n float64, cache *cache.Cache, config *cache.Config, c chan<- types.WebSocketEvent) {
	conf, err := config.GetConfig(context.Background())
	if err != nil {
		return
	}

	arr, err := cache.AddIndexToEcg(context.TODO(), n)
	if err != nil {
		lib.Print(lib.UDP_SERVICE, err)
		return
	}
	length := len(*arr)

	if length%config.ChunkSize > 0 {
		return
	}
	if length < config.ChunkSize*2 {
		return
	}

	if conf.SpectrumUpdateRequest == 1 {
		newConf := *conf
		newConf.SpectrumUpdateRequest = 0
		UpdateSpectrum(&types.EcgSignal{
			Signal: dsp.Signal{
				SampleRate: float64(lib.ECG_HZ),
				Signal:     *arr,
			},
			ChunksSize:         conf.ChunksSize,
			MinPass:            conf.MinPass,
			MaxPass:            conf.MaxPass,
			FilterType:         conf.FilterType,
			WaitSpectrumUpdate: conf.SpectrumUpdateRequest,
		}, c)
		err = config.Set(context.Background(), "config", newConf)
		if err != nil {
			return
		}
		return
	}

	SendHeartBeatData(&types.EcgSignal{
		Signal: dsp.Signal{
			SampleRate: float64(lib.ECG_HZ),
			Signal:     *arr,
		},
		ChunksSize:         conf.ChunksSize,
		MinPass:            conf.MinPass,
		MaxPass:            conf.MaxPass,
		FilterType:         conf.FilterType,
		WaitSpectrumUpdate: conf.SpectrumUpdateRequest,
	}, c)
}
