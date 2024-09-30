package models

import (
	"encoding/json"
	"sia/backend/models/types"
	"time"
)

type (
	Connection struct {
		Uuid       [8]byte   `json:"uuid"`
		CreatedAt  time.Time `json:"created_at"`
		Type types.DataType `json:"type"`
		Data []byte `json:"data"`
	}
)

func (b Connection) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Connection) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &b)
}
