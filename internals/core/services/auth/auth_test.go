package auth

import (
	"os"
	"testing"
)

func setupEnv() {
	os.Setenv("ADMIN_USERNAME", "admin")
	os.Setenv("ADMIN_PASSWORD", "password")
	os.Setenv("JWT_SECRET", "mysecretkey")
}
func TestLogin_Success(t *testing.T) {
	setupEnv()
	authService := NewAuthService()

	token, err := authService.Login("admin", "password")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if token == "" {
		t.Error("Expected a token, got empty string")
	}
}

func TestLogin_Failure(t *testing.T) {
	setupEnv()
	authService := NewAuthService()

	_, err := authService.Login("wronguser", "wrongpass")
	if err == nil {
		t.Error("Expected error for invalid credentials, got nil")
	}
}

func TestValidateToken_Success(t *testing.T) {
	setupEnv()
	authService := NewAuthService()

	token, _ := authService.Login("admin", "password")
	admin, err := authService.ValidateToken(token)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if admin == nil || admin.Username != "admin" {
		t.Error("Expected valid admin, got nil or incorrect data")
	}
}

func TestValidateToken_Failure(t *testing.T) {
	setupEnv()
	authService := NewAuthService()

	invalidToken := "invalid.token.string"
	_, err := authService.ValidateToken(invalidToken)

	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}
