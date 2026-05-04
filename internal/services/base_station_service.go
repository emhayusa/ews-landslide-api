package services

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/models"
	"big-devops-api/internal/repositories"
	"strconv"
)

type BaseStationService interface {
	GetAll(userID uint, role string) ([]dto.BaseStationResponse, error)
	GetByID(id string) (dto.BaseStationResponse, error)
	Create(req *dto.BaseStationRequest) (dto.BaseStationResponse, error)
	Update(id string, req *dto.BaseStationRequest) (dto.BaseStationResponse, error)
	Delete(id string) error
}

type baseStationService struct {
	repo repositories.BaseStationRepository
}

func NewBaseStationService(repo repositories.BaseStationRepository) BaseStationService {
	return &baseStationService{repo}
}

func (s *baseStationService) GetAll(userID uint, role string) ([]dto.BaseStationResponse, error) {
	var stations []models.BaseStation
	var err error

	if role == "mitra" {
		stations, err = s.repo.FindByUser(userID)
	} else {
		stations, err = s.repo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	res := []dto.BaseStationResponse{}
	for _, st := range stations {
		res = append(res, s.mapToResponse(st))
	}
	return res, nil
}

func (s *baseStationService) GetByID(idStr string) (dto.BaseStationResponse, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// Fallback to UUID search if ID is not a number
		st, err := s.repo.FindByUUID(idStr)
		if err != nil {
			return dto.BaseStationResponse{}, err
		}
		return s.mapToResponse(st), nil
	}

	st, err := s.repo.FindByID(uint(id))
	if err != nil {
		return dto.BaseStationResponse{}, err
	}
	return s.mapToResponse(st), nil
}

func (s *baseStationService) Create(req *dto.BaseStationRequest) (dto.BaseStationResponse, error) {
	station := &models.BaseStation{
		Kode:   req.Kode,
		Nama:   req.Nama,
		Lokasi: req.Lokasi,
		Long:   req.Long,
		Lat:    req.Lat,
		SiteID: req.SiteID,
	}

	if err := s.repo.Create(station); err != nil {
		return dto.BaseStationResponse{}, err
	}

	return s.mapToResponse(*station), nil
}

func (s *baseStationService) Update(idStr string, req *dto.BaseStationRequest) (dto.BaseStationResponse, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return dto.BaseStationResponse{}, err
	}

	station, err := s.repo.FindByID(uint(id))
	if err != nil {
		return dto.BaseStationResponse{}, err
	}

	station.Kode = req.Kode
	station.Nama = req.Nama
	station.Lokasi = req.Lokasi
	station.Long = req.Long
	station.Lat = req.Lat
	station.SiteID = req.SiteID

	if err := s.repo.Update(&station); err != nil {
		return dto.BaseStationResponse{}, err
	}

	return s.mapToResponse(station), nil
}

func (s *baseStationService) Delete(idStr string) error {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return err
	}
	return s.repo.Delete(uint(id))
}

func (s *baseStationService) mapToResponse(st models.BaseStation) dto.BaseStationResponse {
	return dto.BaseStationResponse{
		ID:        st.ID,
		UUID:      st.UUID.String(),
		Kode:      st.Kode,
		Nama:      st.Nama,
		Lokasi:    st.Lokasi,
		Long:      st.Long,
		Lat:       st.Lat,
		SiteID:    st.SiteID,
		CreatedAt: st.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: st.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
