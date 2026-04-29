package services

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/models"
	"big-devops-api/internal/repositories"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAll() ([]dto.UserResponse, error)
	GetByUsername(username string) (dto.UserResponse, error)
	Create(req *dto.UserRequest) (dto.UserResponse, error)
	Update(username string, req *dto.UserRequest) (dto.UserResponse, error)
	Delete(username string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetAll() ([]dto.UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []dto.UserResponse
	for _, u := range users {
		res = append(res, dto.UserResponse{
			Username:  u.Username,
			Email:     u.Email,
			FullName:  u.FullName,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
		})
	}
	return res, nil
}

func (s *userService) GetByUsername(username string) (dto.UserResponse, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (s *userService) Create(req *dto.UserRequest) (dto.UserResponse, error) {
	hashedPassword := req.Password
	if req.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		hashedPassword = string(hashed)
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Role:     req.Role,
	}

	if err := s.repo.Create(user); err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) Update(username string, req *dto.UserRequest) (dto.UserResponse, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user.Username = req.Username
	user.Email = req.Email
	user.FullName = req.FullName
	user.Role = req.Role

	if req.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}

	if err := s.repo.Update(&user); err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) Delete(username string) error {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return err
	}
	return s.repo.Delete(fmt.Sprintf("%d", user.ID))
}
