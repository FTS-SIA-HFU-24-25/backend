package types

import "github.com/eripe970/go-dsp-utils"

type EcgSignal struct {
	dsp.Signal
	ChunksSize int
	MinPass    float64
	MaxPass    float64
	FilterType int // 0: no filter, 1: low pass, 2: high pass, 3: band pass
}
