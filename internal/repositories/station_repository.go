package repositories

import (
	"big-devops-api/internal/models"
	"gorm.io/gorm"
)

type StationRepository interface {
	FindAll() ([]models.Station, error)
	FindByUser(userID uint) ([]models.Station, error)
	FindByID(id string) (models.Station, error)
	Create(station *models.Station) error
	Update(station *models.Station) error
	Delete(id string) error
}

type stationRepository struct {
	db *gorm.DB
}

func NewStationRepository(db *gorm.DB) StationRepository {
	return &stationRepository{db}
}

func (r *stationRepository) FindAll() ([]models.Station, error) {
	var stations []models.Station
	err := r.db.Preload("BaseStation").Find(&stations).Error
	return stations, err
}

func (r *stationRepository) FindByUser(userID uint) ([]models.Station, error) {
	var stations []models.Station
	err := r.db.Joins("JOIN sites ON sites.id = stations.site_id").
		Joins("JOIN user_sites ON user_sites.site_id = sites.id").
		Where("user_sites.user_id = ?", userID).
		Preload("BaseStation").
		Find(&stations).Error
	return stations, err
}

func (r *stationRepository) FindByID(id string) (models.Station, error) {
	var station models.Station
	err := r.db.Preload("BaseStation").First(&station, id).Error
	return station, err
}

func (r *stationRepository) Create(station *models.Station) error {
	return r.db.Create(station).Error
}

func (r *stationRepository) Update(station *models.Station) error {
	return r.db.Save(station).Error
}

func (r *stationRepository) Delete(id string) error {
	return r.db.Delete(&models.Station{}, id).Error
}
