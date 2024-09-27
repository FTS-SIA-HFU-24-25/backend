package models

import (
	"encoding/json"
	"time"
)

type (
	Connection struct {
		Uuid       [8]byte   `json:"uuid"`
		CreatedAt  time.Time `json:"created_at"`
		EcgHeader  byte      `json:"ecg_header"`
		GPSHeader  byte      `json:"gps_header"`
		TempHeader byte      `json:"temp_header"`
	}
)

func (b Connection) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Connection) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &b)
}
