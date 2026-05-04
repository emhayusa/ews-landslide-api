package repositories

import (
	"big-devops-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id string) (models.User, error)
	FindByUsername(username string) (models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Sites").Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id string) (models.User, error) {
	var user models.User
	err := r.db.Preload("Sites").First(&user, id).Error
	return user, err
}

func (r *userRepository) FindByUsername(username string) (models.User, error) {
	var user models.User
	err := r.db.Preload("Sites").Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, id).Error
}
