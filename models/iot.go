package models

type (
	UDPEcgData struct {
		Timestamp int64   `json:"timestamp"`
		Value     float64 `json:"value"`
	}
	UDPGpsData struct {
		Timestamp int64   `json:"timestamp"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	UDPTemperatureData struct {
		Timestamp int64   `json:"timestamp"`
		Value     float64 `json:"value"`
	}
)
