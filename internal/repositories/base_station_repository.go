package repositories

import (
	"big-devops-api/internal/models"
	"gorm.io/gorm"
)

type BaseStationRepository interface {
	FindAll() ([]models.BaseStation, error)
	FindByUser(userID uint) ([]models.BaseStation, error)
	FindByID(id uint) (models.BaseStation, error)
	FindByUUID(uuid string) (models.BaseStation, error)
	Create(baseStation *models.BaseStation) error
	Update(baseStation *models.BaseStation) error
	Delete(id uint) error
}

type baseStationRepository struct {
	db *gorm.DB
}

func NewBaseStationRepository(db *gorm.DB) BaseStationRepository {
	return &baseStationRepository{db}
}

func (r *baseStationRepository) FindAll() ([]models.BaseStation, error) {
	var baseStations []models.BaseStation
	err := r.db.Find(&baseStations).Error
	return baseStations, err
}

func (r *baseStationRepository) FindByUser(userID uint) ([]models.BaseStation, error) {
	var baseStations []models.BaseStation
	err := r.db.Joins("JOIN sites ON sites.id = base_stations.site_id").
		Joins("JOIN user_sites ON user_sites.site_id = sites.id").
		Where("user_sites.user_id = ?", userID).
		Find(&baseStations).Error
	return baseStations, err
}

func (r *baseStationRepository) FindByID(id uint) (models.BaseStation, error) {
	var baseStation models.BaseStation
	err := r.db.First(&baseStation, id).Error
	return baseStation, err
}

func (r *baseStationRepository) FindByUUID(uuid string) (models.BaseStation, error) {
	var baseStation models.BaseStation
	err := r.db.Where("uuid = ?", uuid).First(&baseStation).Error
	return baseStation, err
}

func (r *baseStationRepository) Create(baseStation *models.BaseStation) error {
	return r.db.Create(baseStation).Error
}

func (r *baseStationRepository) Update(baseStation *models.BaseStation) error {
	return r.db.Save(baseStation).Error
}

func (r *baseStationRepository) Delete(id uint) error {
	return r.db.Delete(&models.BaseStation{}, id).Error
}
