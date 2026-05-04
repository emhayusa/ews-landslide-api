package services

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/models"
	"big-devops-api/internal/repositories"
	"fmt"
	"strconv"
)

type SiteService interface {
	GetAll(userID uint, role string) ([]dto.SiteResponse, error)
	GetByID(id string, userID uint, role string) (dto.SiteResponse, error)
	Create(req *dto.SiteRequest) (dto.SiteResponse, error)
	Update(id string, req *dto.SiteRequest) (dto.SiteResponse, error)
	Delete(id string) error
}

type siteService struct {
	repo repositories.SiteRepository
}

func NewSiteService(repo repositories.SiteRepository) SiteService {
	return &siteService{repo}
}

func (s *siteService) GetAll(userID uint, role string) ([]dto.SiteResponse, error) {
	var sites []models.Site
	var err error

	if role == "mitra" {
		sites, err = s.repo.FindByUser(userID)
	} else {
		sites, err = s.repo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	res := []dto.SiteResponse{}
	for _, site := range sites {
		res = append(res, s.mapToResponse(site))
	}
	return res, nil
}

func (s *siteService) GetByID(idStr string, userID uint, role string) (dto.SiteResponse, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return dto.SiteResponse{}, err
	}

	site, err := s.repo.FindByID(uint(id))
	if err != nil {
		return dto.SiteResponse{}, err
	}

	// Access control for mitra
	if role == "mitra" {
		hasAccess := false
		for _, u := range site.Users {
			if u.ID == userID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return dto.SiteResponse{}, fmt.Errorf("access denied")
		}
	}

	return s.mapToResponse(site), nil
}

func (s *siteService) Create(req *dto.SiteRequest) (dto.SiteResponse, error) {
	site := &models.Site{
		Nama:       req.Nama,
		Lokasi:     req.Lokasi,
		Keterangan: req.Keterangan,
	}

	if err := s.repo.Create(site); err != nil {
		return dto.SiteResponse{}, err
	}

	return s.mapToResponse(*site), nil
}

func (s *siteService) Update(idStr string, req *dto.SiteRequest) (dto.SiteResponse, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return dto.SiteResponse{}, err
	}

	site, err := s.repo.FindByID(uint(id))
	if err != nil {
		return dto.SiteResponse{}, err
	}

	site.Nama = req.Nama
	site.Lokasi = req.Lokasi
	site.Keterangan = req.Keterangan

	if err := s.repo.Update(&site); err != nil {
		return dto.SiteResponse{}, err
	}

	return s.mapToResponse(site), nil
}

func (s *siteService) Delete(idStr string) error {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return err
	}
	return s.repo.Delete(uint(id))
}

func (s *siteService) mapToResponse(site models.Site) dto.SiteResponse {
	baseStations := []dto.BaseStationResponse{}
	for _, bs := range site.BaseStations {
		baseStations = append(baseStations, dto.BaseStationResponse{
			ID:        bs.ID,
			UUID:      bs.UUID.String(),
			Kode:      bs.Kode,
			Nama:      bs.Nama,
			Lokasi:    bs.Lokasi,
			Long:      bs.Long,
			Lat:       bs.Lat,
			CreatedAt: bs.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: bs.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	stations := []dto.StationResponse{}
	for _, st := range site.Stations {
		stations = append(stations, dto.StationResponse{
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
		})
	}

	return dto.SiteResponse{
		ID:           site.ID,
		Nama:         site.Nama,
		Lokasi:       site.Lokasi,
		Keterangan:   site.Keterangan,
		BaseStations: baseStations,
		Stations:     stations,
		CreatedAt:    site.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    site.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
