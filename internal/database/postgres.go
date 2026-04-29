package database

import (
	"big-devops-api/internal/config"
	"big-devops-api/internal/models"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("Connected to PostgreSQL")

	// Auto migration
	err = db.AutoMigrate(&models.User{}, &models.Station{}, &models.Monitoring{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// Seed Admin User if none exist
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("presisi13!@#"), bcrypt.DefaultCost)
		admin := models.User{
			Username: "admin",
			Email:    "admin@presisipedia.xyz",
			Password: string(hashedPassword),
			FullName: "Admin",
			Role:     "admin",
		}
		db.Create(&admin)
		fmt.Println("Seeded default admin user: kelvin / password123")
	}

	DB = db
}
