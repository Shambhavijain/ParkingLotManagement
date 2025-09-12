package auth

import (
	"errors"
	"os"
	"testing"
	"time"

	"parkingSlotManagement/internals/core/domain"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository mocks the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByUsername(username string) (*domain.Admin, error) {
	args := m.Called(username)
	if admin, ok := args.Get(0).(*domain.Admin); ok {
		return admin, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestLogin_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)

	admin := &domain.Admin{
		ID:       "admin-1",
		Username: "admin",
		Password: "admin123",
	}

	mockRepo.On("GetByUsername", "admin").Return(admin, nil)

	token, err := authService.Login("admin", "admin123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token format
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("testsecret"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)

	admin := &domain.Admin{
		ID:       "admin-1",
		Username: "admin",
		Password: "admin123",
	}

	mockRepo.On("GetByUsername", "admin").Return(admin, nil)

	token, err := authService.Login("admin", "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)

	mockRepo.On("GetByUsername", "admin").Return(nil, errors.New("not found"))

	token, err := authService.Login("admin", "admin123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "user not found", err.Error())
}

func TestValidateToken_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)

	// Create a valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": "admin-1",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte("testsecret"))

	admin, err := authService.ValidateToken(tokenStr)

	assert.NoError(t, err)
	assert.Equal(t, "admin-1", admin.ID)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)

	invalidToken := "this.is.not.valid"

	admin, err := authService.ValidateToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, admin)
	assert.Equal(t, "invalid token", err.Error())
}
