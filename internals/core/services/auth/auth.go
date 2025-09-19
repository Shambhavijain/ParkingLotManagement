package auth

import (
	"errors"
	"os"
	"parkingSlotManagement/internals/core/domain"
	"parkingSlotManagement/internals/ports"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthInterface interface {
	Login(username, password string) (string, error)
	ValidateToken(tokenStr string) (*domain.Admin, error)
}
type AuthService struct {
	repo      ports.UserRepository
	secretKey string
}

func NewAuthService(repo ports.UserRepository) *AuthService {
	return &AuthService{
		repo:      repo,
		secretKey: os.Getenv("JWT_SECRET"),
	}
}

func (a *AuthService) Login(username, password string) (string, error) {
	admin, err := a.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if password != admin.Password {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": admin.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(a.secretKey))
}

func (a *AuthService) ValidateToken(tokenStr string) (*domain.Admin, error) {
	tokenStr = strings.TrimSpace(tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(a.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	adminID, ok := claims["admin_id"].(string)
	if !ok {
		return nil, errors.New("invalid admin ID in token")
	}

	return &domain.Admin{ID: adminID}, nil
}
