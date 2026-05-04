package repositories

import (
	"big-devops-api/internal/models"
	"gorm.io/gorm"
)

type DeformationRepository interface {
	Create(deformation *models.Deformation) error
	GetLatestByStation(obsCode string) (models.Deformation, error)
}

type deformationRepository struct {
	db *gorm.DB
}

func NewDeformationRepository(db *gorm.DB) DeformationRepository {
	return &deformationRepository{db}
}

func (r *deformationRepository) Create(deformation *models.Deformation) error {
	return r.db.Create(deformation).Error
}

func (r *deformationRepository) GetLatestByStation(obsCode string) (models.Deformation, error) {
	var deformation models.Deformation
	err := r.db.Where("obs_code = ?", obsCode).Order("ts desc").First(&deformation).Error
	return deformation, err
}
