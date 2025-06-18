package types

const (
	// Read more -> "/translator/udp.go"
	UDP_EKG_SENSOR         = 0
	UDP_TEMPERATURE_SENSOR = 1
	UDP_GYRO_SENSOR        = 2
	UDP_ACCEL_SENSOR        = 3 
	END_CONNECTION         = 4 
)

type (
	EKG_SENSOR struct {
		Value     float64   `json:"value"`
	}
	TEMPERATURE_SENSOR struct {
		Value     float32   `json:"value"`
	}
	GYRO_SENSOR struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
		Z float32 `json:"z"`
	}
	ACCEL_SENSOR struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
		Z float32 `json:"z"`
	}
	END_REQUEST struct {
		SENSOR_ID int `json:"sensor_id"`
	}
)

func (s *EKG_SENSOR) Init(value float64) *EKG_SENSOR {
	s.Value = value
	return s
}

func (s *TEMPERATURE_SENSOR) Init(value float32) *TEMPERATURE_SENSOR {
	s.Value = value
	return s
}

func (s *GYRO_SENSOR) Init(x, y, z float32) *GYRO_SENSOR {
	s.X = x
	s.Y = y
	s.Z = z
	return s
}

func (s *ACCEL_SENSOR) Init(x, y, z float32) *ACCEL_SENSOR {
	s.X = x
	s.Y = y
	s.Z = z
	return s
}

