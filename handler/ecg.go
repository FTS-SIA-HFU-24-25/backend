package handler

import (
	"sia/backend/types"
	"time"

	"github.com/eripe970/go-dsp-utils"
)

type HeartBeat struct {
	Count     int       `json:"count"`
	Avg       int       `json:"avg"`
	Timestamp time.Time `json:"timestamp"`
}

func InitECGHandler(ecgChan <-chan dsp.Signal, wsChan chan<- types.WebSocketEvent) {
	for c := range ecgChan {
		if c.Length() < int(c.SampleRate*5) {
			continue
		}
		if c.Length()%int(c.SampleRate*5) > 0 {
			continue
		}
		go updateHeartBeat(getLastSeconds(&c, time.Second*5), wsChan)
	}
}

func updateHeartBeat(s *dsp.Signal, wsChan chan<- types.WebSocketEvent) {
	rPeak := dsp.GetRPeaks(s)
	events := types.WebSocketEvent{
		Event: "heartbeat",
		Data: &HeartBeat{
			Count:     rPeak.Count(),
			Avg:       rPeak.Avg(),
			Timestamp: time.Now(),
		},
	}
	wsChan <- events
}

func getLastSeconds(s *dsp.Signal, dur time.Duration) *dsp.Signal {
	sample := make([]float64, 0)
	for i := s.Length() - int(s.SampleRate*dur.Seconds()); i < s.Length(); i++ {
		sample = append(sample, s.Signal[i])
	}

	return &dsp.Signal{
		SampleRate: s.SampleRate,
		Signal:     sample,
	}
}
