package types

import "github.com/eripe970/go-dsp-utils"

type EcgSignal struct {
	dsp.Signal
	ChunksSize int
}
