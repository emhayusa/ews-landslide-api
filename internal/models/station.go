package models

import (
	"time"

	"gorm.io/gorm"
)

type Station struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	StationID    string         `gorm:"uniqueIndex;not null" json:"station_id"` // Code based unique ID
	Name         string         `gorm:"not null" json:"name"`
	Latitude     float64        `json:"latitude"`
	Longitude    float64        `json:"longitude"`
	Location     string         `json:"location"`
	Description  string         `json:"description"`
	Status       string         `json:"status"` // active, inactive, maintenance
	HardwareType string         `json:"hardware_type"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
