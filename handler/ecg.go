package handler

import (
	"sia/backend/lib"
	"sia/backend/types"
	"slices"
	"time"

	"github.com/eripe970/go-dsp-utils"
)

type (
	HeartBeat struct {
		Count     int       `json:"count"`
		Avg       int       `json:"avg"`
		Timestamp time.Time `json:"timestamp"`
	}
	EcgWSEvent struct {
		Signals   []float64 `json:"signals"`
		Avg       int       `json:"avg"`
		Max       float64   `json:"max"`
		Min       float64   `json:"min"`
		Frequency int       `json:"frequency"`
		Timestamp time.Time `json:"timestamp"`
	}
)

func InitECGHandler(ecgChan <-chan dsp.Signal, wsChan chan<- types.WebSocketEvent) {
	for c := range ecgChan {
		if c.Length() > int(c.SampleRate)*2 && c.Length()%(int(c.SampleRate)/10) == 0 {
			sendHeartBeatData(&c, wsChan)
		}
	}
}

func sendHeartBeatData(c *dsp.Signal, wsChan chan<- types.WebSocketEvent) {
	lib.Print(lib.ECG_SERVICE, c.SampleRate, c.Signal)
	length := int(c.SampleRate) / 10

	data, err := c.Normalize()
	if err != nil {
		return
	}

	spectrum, _ := data.FrequencySpectrum()
	maxV := spectrum.Max()
	// minV := spectrum.Min()
	if maxV == 0 {
		return
	}

	//data, _ = c.HighPassFilter(maxV * 2)
	// data, _ = c.BandPassFilter(0.5, 5)

	rPeak := dsp.GetRPeaks(c)

	wsChan <- types.WebSocketEvent{
		Event: "ekg-changes",
		Data: EcgWSEvent{
			Signals:   data.Signal[len(data.Signal)-length:],
			Avg:       rPeak.Avg(),
			Max:       maxV * 2,
			Min:       maxV * -2,
			Frequency: int(data.SampleRate),
			Timestamp: time.Now(),
		},
	}
}

func calculateSpectrumPeakFreq(d []float64, maxV float64) float64 {
	lib.Print(lib.ECG_SERVICE, len(d), int(maxV))
	seq := float64(len(d)) / maxV

	searchData := d[int(seq)*2:]

	maxVal := slices.Max(searchData)

	return maxVal
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
