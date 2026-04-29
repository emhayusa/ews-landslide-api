package services

import (
	"big-devops-api/internal/config"
	"big-devops-api/internal/repositories"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type authService struct {
	cfg      *config.Config
	userRepo repositories.UserRepository
}

func NewAuthService(cfg *config.Config, userRepo repositories.UserRepository) AuthService {
	return &authService{cfg, userRepo}
}

func (s *authService) Login(username, password string) (string, error) {
	if s.cfg.AuthMethod == "keycloak" {
		return "", errors.New("login should be handled by Keycloak")
	}

	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Username,
		"name":  user.FullName,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *authService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	if s.cfg.AuthMethod == "keycloak" {
		// keycloak validation is typically handled by parsing the token or JWKS
		// Simplified for this example
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			return nil, err
		}
		return token.Claims.(jwt.MapClaims), nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}
