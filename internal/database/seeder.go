package database

import (
	"big-devops-api/internal/models"
	"fmt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	// 1. Seed Site
	var site models.Site
	if err := db.Where("nama = ?", "Wilayah Semarang").First(&site).Error; err != nil {
		site = models.Site{
			Nama:   "Wilayah Semarang",
			Lokasi: "Semarang, Jawa Tengah",
		}
		db.Create(&site)
		fmt.Println("Seeded Site: Wilayah Semarang")
	}

	// 2. Seed Base Station
	var base models.BaseStation
	if err := db.Where("kode = ?", "CSEM").First(&base).Error; err != nil {
		base = models.BaseStation{
			Kode:   "CSEM",
			Nama:   "Base Semarang",
			Lokasi: "Semarang",
			Lat:    -7.0500,
			Long:   110.4300,
			SiteID: &site.ID,
		}
		db.Create(&base)
		fmt.Println("Seeded Base Station: CSEM")
	}

	// 3. Seed Station (Rover)
	var rover models.Station
	if err := db.Where("station_id = ?", "UNGR").First(&rover).Error; err != nil {
		rover = models.Station{
			StationID:       "UNGR",
			Name:            "Rover Ungaran",
			Latitude:        -7.1200,
			Longitude:       110.4000,
			Location:        "Ungaran",
			Status:          "ACTIVE",
			BaseStationID:   &base.ID,
			SiteID:          &site.ID,
			InitialDistance: 17528.774,
			URLStreaming:    "ws://36.92.41.75:8000/ws/data?ref=CSEM&obs=UNGR",
		}
		db.Create(&rover)
		fmt.Println("Seeded Station: UNGR")
	}
}
