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
		Signals []float64 `json:"signals"`
		Avg     int       `json:"avg"`
	}
	SpectrumWSEvent struct {
		Spectrum  []float64 `json:"spectrum"`
		Frequency []float64 `json:"frequency"`
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

	if len(data.Signal)%int(data.SampleRate*10) == 0 {
		newData := getLastSeconds(data, time.Second*8)
		spectrum, _ := newData.FrequencySpectrum()
		wsChan <- types.WebSocketEvent{
			Event: "spectrum-changes",
			Data: SpectrumWSEvent{
				Spectrum:  spectrum.Spectrum,
				Frequency: spectrum.Frequencies,
			},
		}
	}

	if c.FilterType != 0 {
		switch c.FilterType {
		case 1:
			data, err = c.LowPassFilter(c.MaxPass)
		case 2:
			data, err = c.HighPassFilter(c.MinPass)
		case 3:
			data, err = c.BandPassFilter(c.MinPass, c.MaxPass)
		}
		if err != nil {
			return
		}
	}

	rPeak := dsp.GetRPeaks(&c.Signal)

	wsChan <- types.WebSocketEvent{
		Event: "ekg-changes",
		Data: EcgWSEvent{
			Signals: data.Signal[len(data.Signal)-length:],
			Avg:     rPeak.Avg(),
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
	length := int(s.SampleRate * dur.Seconds())
	if length >= len(s.Signal) {
		return s
	}
	s.Signal = s.Signal[len(s.Signal)-length:]
	return s
}
