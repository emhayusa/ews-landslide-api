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
	Status          string         `json:"status"` // active, inactive, maintenance
	BaseStationID   *uint          `json:"base_station_id"`
	BaseStation     *BaseStation   `gorm:"foreignKey:BaseStationID" json:"base_station,omitempty"`
	InitialDistance float64        `json:"initial_distance"`
	URLStreaming    string         `json:"url_streaming"`
	SiteID          *uint          `gorm:"column:site_id" json:"site_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
