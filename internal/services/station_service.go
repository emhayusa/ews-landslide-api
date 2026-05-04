package services

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/models"
	"big-devops-api/internal/repositories"
)

type StationService interface {
	GetAll(userID uint, role string) ([]dto.StationResponse, error)
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

func (s *stationService) GetAll(userID uint, role string) ([]dto.StationResponse, error) {
	var stations []models.Station
	var err error

	if role == "mitra" {
		stations, err = s.repo.FindByUser(userID)
	} else {
		stations, err = s.repo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	res := []dto.StationResponse{}
	for _, st := range stations {
		res = append(res, s.mapToResponse(st))
	}
	return res, nil
}

func (s *stationService) GetByID(id string) (dto.StationResponse, error) {
	st, err := s.repo.FindByID(id)
	if err != nil {
		return dto.StationResponse{}, err
	}
	return s.mapToResponse(st), nil
}

func (s *stationService) Create(req *dto.StationRequest) (dto.StationResponse, error) {
	station := &models.Station{
		StationID:    req.StationID,
		Name:         req.Name,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Location:     req.Location,
		Description:  req.Description,
		Status:          req.Status,
		BaseStationID:   req.BaseStationID,
		InitialDistance: req.InitialDistance,
		URLStreaming:    req.URLStreaming,
		SiteID:          req.SiteID,
	}

	if err := s.repo.Create(station); err != nil {
		return dto.StationResponse{}, err
	}

	return s.mapToResponse(*station), nil
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
	station.BaseStationID = req.BaseStationID
	station.InitialDistance = req.InitialDistance
	station.URLStreaming = req.URLStreaming
	station.SiteID = req.SiteID

	if err := s.repo.Update(&station); err != nil {
		return dto.StationResponse{}, err
	}

	return s.mapToResponse(station), nil
}

func (s *stationService) mapToResponse(st models.Station) dto.StationResponse {
	res := dto.StationResponse{
		ID:              st.ID,
		StationID:       st.StationID,
		Name:            st.Name,
		Latitude:        st.Latitude,
		Longitude:       st.Longitude,
		Location:        st.Location,
		Status:          st.Status,
		BaseStationID:   st.BaseStationID,
		InitialDistance: st.InitialDistance,
		URLStreaming:    st.URLStreaming,
		SiteID:          st.SiteID,
	}

	if st.BaseStation != nil {
		res.BaseStation = &dto.BaseStationResponse{
			ID:     st.BaseStation.ID,
			UUID:   st.BaseStation.UUID.String(),
			Kode:   st.BaseStation.Kode,
			Nama:   st.BaseStation.Nama,
			Lokasi: st.BaseStation.Lokasi,
			Long:   st.BaseStation.Long,
			Lat:    st.BaseStation.Lat,
		}
	}

	return res
}

func (s *stationService) Delete(id string) error {
	return s.repo.Delete(id)
}
