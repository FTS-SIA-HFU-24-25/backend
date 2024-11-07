package types

import "time"

const (
	// Read more -> "/translator/udp.go"
	UDP_EKG_SENSOR         = 0
	UDP_TEMPERATURE_SENSOR = 1
	UDP_GPS_SERVICE        = 2
	END_CONNECTION         = 3
)

type (
	EKG_SENSOR struct {
		Value     float64   `json:"value"`
		Timestamp time.Time `json:"timestamp"`
	}
	TEMPERATURE_SENSOR struct {
		Value     float64   `json:"value"`
		Timestamp time.Time `json:"timestamp"`
	}
	GPS_SERVICE struct {
		Latitude  float64   `json:"latitude"`
		Longitude float64   `json:"longitude"`
		Timestamp time.Time `json:"timestamp"`
	}
	END_REQUEST struct {
		SENSOR_ID int `json:"sensor_id"`
	}
)

func (s *EKG_SENSOR) Init(value float64) *EKG_SENSOR {
	s.Value = value
	s.Timestamp = time.Now()
	return s
}

func (s *TEMPERATURE_SENSOR) Init(value float64) *TEMPERATURE_SENSOR {
	s.Value = value
	s.Timestamp = time.Now()
	return s
}

func (s *GPS_SERVICE) Init(latitude, longitude float64) *GPS_SERVICE {
	s.Latitude = latitude
	s.Longitude = longitude
	s.Timestamp = time.Now()
	return s
}
