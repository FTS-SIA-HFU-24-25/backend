package lib

import (
	"bytes"
	"encoding/binary"
	"math"
)

func Float64FromBytes(b []byte) float64 {
	bits := binary.LittleEndian.Uint64(b)
	float := math.Float64frombits(bits)
	return float
}

func Float64ToBytes(f float64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, f)
	if err != nil {
		Print(UDP_SERVICE, "Failed to convert float64 to bytes")
		return nil
	}
	return buf.Bytes()
}
