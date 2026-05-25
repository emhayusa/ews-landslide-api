package services

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/repositories"
)

type MonitoringService interface {
	GetRainfallHistory(stationID string, limit int) ([]dto.RainfallHistoryResponse, error)
}

type monitoringService struct {
	repo repositories.MonitoringRepository
}

func NewMonitoringService(repo repositories.MonitoringRepository) MonitoringService {
	return &monitoringService{repo}
}

func (s *monitoringService) GetRainfallHistory(stationID string, limit int) ([]dto.RainfallHistoryResponse, error) {
	monitorings, err := s.repo.FindByStationID(stationID, limit)
	if err != nil {
		return nil, err
	}

	res := make([]dto.RainfallHistoryResponse, 0)
	for _, m := range monitorings {
		res = append(res, dto.RainfallHistoryResponse{
			Timestamp: m.Timestamp,
			Hourly:    m.CurahHujanHourly,
			Daily:     m.CurahHujanDaily,
		})
	}
	return res, nil
}
