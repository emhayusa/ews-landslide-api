package repositories

import (
	"big-devops-api/internal/models"
	"gorm.io/gorm"
)

type MonitoringRepository interface {
	Create(monitoring *models.Monitoring) error
}

type monitoringRepository struct {
	db *gorm.DB
}

func NewMonitoringRepository(db *gorm.DB) MonitoringRepository {
	return &monitoringRepository{db: db}
}

func (r *monitoringRepository) Create(monitoring *models.Monitoring) error {
	return r.db.Create(monitoring).Error
}
