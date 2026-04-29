package services

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/models"
	"big-devops-api/internal/repositories"
)

type StationService interface {
	GetAll() ([]dto.StationResponse, error)
	GetByID(id string) (dto.StationResponse, error)
	Create(req *dto.StationRequest) (dto.StationResponse, error)
	Update(id string, req *dto.StationRequest) (dto.StationResponse, error)
	Delete(id string) error
}

type stationService struct {
	repo repositories.StationRepository
}

func NewStationService(repo repositories.StationRepository) StationService {
	return &stationService{repo}
}

func (s *stationService) GetAll() ([]dto.StationResponse, error) {
	stations, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	res := []dto.StationResponse{}
	for _, st := range stations {
		res = append(res, dto.StationResponse{
			ID:           st.ID,
			StationID:    st.StationID,
			Name:         st.Name,
			Latitude:     st.Latitude,
			Longitude:    st.Longitude,
			Location:     st.Location,
			Status:       st.Status,
			HardwareType: st.HardwareType,
		})
	}
	return res, nil
}

func (s *stationService) GetByID(id string) (dto.StationResponse, error) {
	st, err := s.repo.FindByID(id)
	if err != nil {
		return dto.StationResponse{}, err
	}
	return dto.StationResponse{
		ID:           st.ID,
		StationID:    st.StationID,
		Name:         st.Name,
		Latitude:     st.Latitude,
		Longitude:    st.Longitude,
		Location:     st.Location,
		Status:       st.Status,
		HardwareType: st.HardwareType,
	}, nil
}

func (s *stationService) Create(req *dto.StationRequest) (dto.StationResponse, error) {
	station := &models.Station{
		StationID:    req.StationID,
		Name:         req.Name,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Location:     req.Location,
		Description:  req.Description,
		Status:       req.Status,
		HardwareType: req.HardwareType,
	}

	if err := s.repo.Create(station); err != nil {
		return dto.StationResponse{}, err
	}

	return dto.StationResponse{
		ID:           station.ID,
		StationID:    station.StationID,
		Name:         station.Name,
		Latitude:     station.Latitude,
		Longitude:    station.Longitude,
		Location:     station.Location,
		Status:       station.Status,
		HardwareType: station.HardwareType,
	}, nil
}

func (s *stationService) Update(id string, req *dto.StationRequest) (dto.StationResponse, error) {
	station, err := s.repo.FindByID(id)
	if err != nil {
		return dto.StationResponse{}, err
	}

	station.StationID = req.StationID
	station.Name = req.Name
	station.Latitude = req.Latitude
	station.Longitude = req.Longitude
	station.Location = req.Location
	station.Description = req.Description
	station.Status = req.Status
	station.HardwareType = req.HardwareType

	if err := s.repo.Update(&station); err != nil {
		return dto.StationResponse{}, err
	}

	return dto.StationResponse{
		ID:           station.ID,
		StationID:    station.StationID,
		Name:         station.Name,
		Latitude:     station.Latitude,
		Longitude:    station.Longitude,
		Location:     station.Location,
		Status:       station.Status,
		HardwareType: station.HardwareType,
	}, nil
}

func (s *stationService) Delete(id string) error {
	return s.repo.Delete(id)
}
