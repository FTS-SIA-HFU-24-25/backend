package handler

import (
	"context"
	"log"
	"sia/backend/cache"
	"sia/backend/lib"
	"sia/backend/types"

	"github.com/eripe970/go-dsp-utils"
)

type (
	EcgWSEvent struct {
		Signals []float64 `json:"signals"`
		Avg     int       `json:"avg"`
	}
	SpectrumWSEvent struct {
		Spectrum  []float64 `json:"spectrum"`
		Frequency []float64 `json:"frequency"`
	}
)

func updateEcgChannel(n float64, cache *cache.Cache, config *cache.Config, c chan<- types.WebSocketEvent) {
	ctx := context.Background()
	conf, err := config.GetConfig(ctx)
	if err != nil {
		log.Printf("[ERROR] Failed to get config: %v", err)
		return
	}

	// Append to buffer
	signalBuf, err := cache.AddIndexToEcg(ctx, n)
	if err != nil {
		go log.Printf("[UDP_SERVICE] Cache error: %v", err)
		return
	}

	length := len(signalBuf)
	if length < conf.ChunksSize*2 || length%conf.ChunksSize != 0 {
		return
	}

	// Trim buffer to prevent unbounded growth
	if length > 10000 {
		signalBuf = signalBuf[length-10000:]
	}

	signal := &types.EcgSignal{
		Signal: dsp.Signal{
			SampleRate: float64(lib.ECG_HZ),
			Signal:     signalBuf,
		},
		ChunksSize: conf.ChunksSize,
	}

	// Handle heartbeat data
	processHeartBeat(signal, c)
}

func processHeartBeat(c *types.EcgSignal, wsChan chan<- types.WebSocketEvent) {
	data, err := c.Normalize()
	if err != nil {
		log.Printf("[ERROR] Normalize failed: %v", err)
		return
	}

	rPeak := dsp.GetRPeaks(&c.Signal)
	if c.ChunksSize <= 0 || c.ChunksSize > len(data.Signal) {
		log.Printf("[ERROR] Invalid chunk size: %d, signal length: %d", c.ChunksSize, len(data.Signal))
		return
	}

	select {
	case wsChan <- types.WebSocketEvent{
		Event: "ekg-changes",
		Data: EcgWSEvent{
			Signals: data.Signal[len(data.Signal)-c.ChunksSize:],
			Avg:     rPeak.Avg(),
		},
	}:
	default:
		log.Println("[ERROR] Dropped ekg event: channel full")
	}
}
