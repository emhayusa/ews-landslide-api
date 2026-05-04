package models

import (
	"time"
)

type Monitoring struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	StationReferID   uint      `gorm:"column:station_id;not null" json:"station_id"`
	Station          Station   `gorm:"foreignKey:StationReferID;references:ID" json:"-"`
	Timestamp        time.Time `json:"timestamp"`
	CurahHujanDaily  float64   `json:"curah_hujan_daily"`
	Baterai          float64   `json:"baterai"`
	Solar            float64   `json:"solar"`
	Alarm            int       `json:"alarm"`
	CurahHujanHourly float64   `json:"curah_hujan_hourly"`
	RawPayload       string    `gorm:"type:text" json:"-"`
	CreatedAt        time.Time `json:"created_at"`
}
