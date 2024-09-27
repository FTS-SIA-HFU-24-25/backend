package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        [8]byte        `json:"id" gorm:"primary_key"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Hardware struct {
	Model
	Ecg         []EcgData         `json:"ecg" gorm:"foreignKey:HardwareID"`
	GPS         []GpsData         `json:"gps" gorm:"foreignKey:HardwareID"`
	Temperature []TemperatureData `json:"temperature" gorm:"foreignKey:HardwareID"`
}

type EcgData struct {
	Model
	HardwareID [8]byte `json:"hardware_id"`
	Timestamp  int64   `json:"timestamp"`
	Value      float64 `json:"value"`
}

type TemperatureData struct {
	Model
	HardwareID [8]byte `json:"hardware_id"`
	Timestamp  int64   `json:"timestamp"`
	Value      int     `json:"value"`
}

type GpsData struct {
	Model
	HardwareID [8]byte `json:"hardware_id"`
	Timestamp  int64   `json:"timestamp"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}
