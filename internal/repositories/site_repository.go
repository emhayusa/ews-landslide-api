package repositories

import (
	"big-devops-api/internal/models"
	"gorm.io/gorm"
)

type SiteRepository interface {
	FindAll() ([]models.Site, error)
	FindByUser(userID uint) ([]models.Site, error)
	FindByID(id uint) (models.Site, error)
	Create(site *models.Site) error
	Update(site *models.Site) error
	Delete(id uint) error
}

type siteRepository struct {
	db *gorm.DB
}

func NewSiteRepository(db *gorm.DB) SiteRepository {
	return &siteRepository{db}
}

func (r *siteRepository) FindAll() ([]models.Site, error) {
	var sites []models.Site
	err := r.db.Preload("BaseStations").Preload("Stations").Find(&sites).Error
	return sites, err
}

func (r *siteRepository) FindByUser(userID uint) ([]models.Site, error) {
	var sites []models.Site
	err := r.db.Joins("JOIN user_sites ON user_sites.site_id = sites.id").
		Where("user_sites.user_id = ?", userID).
		Preload("BaseStations").Preload("Stations").
		Find(&sites).Error
	return sites, err
}

func (r *siteRepository) FindByID(id uint) (models.Site, error) {
	var site models.Site
	err := r.db.Preload("BaseStations").Preload("Stations").First(&site, id).Error
	return site, err
}

func (r *siteRepository) Create(site *models.Site) error {
	return r.db.Create(site).Error
}

func (r *siteRepository) Update(site *models.Site) error {
	return r.db.Save(site).Error
}

func (r *siteRepository) Delete(id uint) error {
	return r.db.Delete(&models.Site{}, id).Error
}
