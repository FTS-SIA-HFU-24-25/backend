package handler

import (
	"sia/backend/lib"
	"sia/backend/types"
	"slices"
	"time"

	"github.com/eripe970/go-dsp-utils"
)

type (
	EcgWSEvent struct {
		Signals   []float64 `json:"signals"`
		Avg       int       `json:"avg"`
		Max       float64   `json:"max"`
		Min       float64   `json:"min"`
		Frequency int       `json:"frequency"`
		Timestamp time.Time `json:"timestamp"`
	}
)

func InitECGHandler(ecgChan <-chan types.EcgSignal, wsChan chan<- types.WebSocketEvent) {
	for c := range ecgChan {
		sendHeartBeatData(&c, wsChan)
	}
}

func sendHeartBeatData(c *types.EcgSignal, wsChan chan<- types.WebSocketEvent) {
	lib.Print(lib.ECG_SERVICE, c.SampleRate, c.Signal)
	length := int(c.SampleRate) / c.ChunksSize

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

	rPeak := dsp.GetRPeaks(&c.Signal)

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
