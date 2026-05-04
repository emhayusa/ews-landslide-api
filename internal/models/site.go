package models

import (
	"time"

	"gorm.io/gorm"
)

type Site struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Nama         string         `gorm:"not null" json:"nama"`
	Lokasi       string         `json:"lokasi"`
	Keterangan   string         `json:"keterangan"`
	BaseStations []BaseStation  `gorm:"foreignKey:SiteID" json:"base_stations,omitempty"`
	Stations     []Station      `gorm:"foreignKey:SiteID" json:"stations,omitempty"`
	Users        []User         `gorm:"many2many:user_sites;" json:"users,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
