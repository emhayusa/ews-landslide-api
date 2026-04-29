package models

import (
	"time"
)

type Monitoring struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	StationReferID uint      `gorm:"column:station_id;not null" json:"station_id"`
	Station        Station   `gorm:"foreignKey:StationReferID;references:ID" json:"-"`
	Timestamp   time.Time `json:"timestamp"`
	Bucket      float64   `json:"bucket"`
	Battery     float64   `json:"battery"`
	Solar       float64   `json:"solar"`
	Alarm       int       `json:"alarm"`
	MaxBucket   float64   `json:"max_bucket"`
	Deformasi   float64   `json:"deformasi"`
	RawPayload  string    `gorm:"type:text" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
}
