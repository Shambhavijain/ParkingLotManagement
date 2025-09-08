package auth

import (
	"errors"
	"fmt"
	"os"
	"parkingSlotManagement/internals/core/domain"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(username, password string) (string, error)
	ValidateToken(token string) (*domain.Admin, error)
}

type AuthServiceImpl struct {
	admin     domain.Admin
	secretKey string
}

func NewAuthService() *AuthServiceImpl {
	return &AuthServiceImpl{
		admin: domain.Admin{
			ID:       "admin-1",
			Username: os.Getenv("ADMIN_USERNAME"),
			Password: os.Getenv("ADMIN_PASSWORD"),
		},
		secretKey: os.Getenv("JWT_SECRET"),
	}
}

func (a *AuthServiceImpl) Login(username, password string) (string, error) {
	if username != a.admin.Username || password != a.admin.Password {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": a.admin.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(a.secretKey))
}

func (a *AuthServiceImpl) ValidateToken(tokenStr string) (*domain.Admin, error) {
	tokenStr = strings.TrimSpace(tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})

	if err != nil {
		fmt.Println("Token parse error:", err)
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return nil, errors.New("invalid token")
	}

	return &a.admin, nil
}
