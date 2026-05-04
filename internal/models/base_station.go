package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseStation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UUID      uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	Kode      string         `gorm:"uniqueIndex;not null" json:"kode"`
	Nama      string         `gorm:"not null" json:"nama"`
	Lokasi    string         `json:"lokasi"`
	Long      float64        `json:"long"`
	Lat       float64        `json:"lat"`
	SiteID    *uint          `gorm:"column:site_id" json:"site_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
