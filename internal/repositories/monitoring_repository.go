package repositories

import (
	"big-devops-api/internal/models"
	"gorm.io/gorm"
)

type MonitoringRepository interface {
	Create(monitoring *models.Monitoring) error
	FindByStationID(stationID string, limit int) ([]models.Monitoring, error)
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

func (r *monitoringRepository) FindByStationID(stationID string, limit int) ([]models.Monitoring, error) {
	var monitorings []models.Monitoring
	err := r.db.Where("station_id = (SELECT id FROM stations WHERE station_id = ?)", stationID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&monitorings).Error
	return monitorings, err
}
