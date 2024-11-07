package types

import "github.com/eripe970/go-dsp-utils"

type EcgSignal struct {
	dsp.Signal
}

func (s *EcgSignal) AddValue(v float64) *EcgSignal {
	if s.Length() >= 3000 {
		s.Signal.Signal = s.Signal.Signal[1:]
	}
	s.Signal.Signal = append(s.Signal.Signal, v)

	return s
}
