package models

import (
	"time"
)

type Deformation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TS        time.Time `gorm:"not null" json:"ts"`
	Distance  float64   `json:"distance"`
	Offset    float64   `json:"offset"`
	RefCode   string    `gorm:"index" json:"ref_code"`
	ObsCode   string    `gorm:"index" json:"obs_code"`
	CreatedAt time.Time `json:"created_at"`
}
