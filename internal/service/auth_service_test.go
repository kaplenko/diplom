package service

import (
	"testing"

	"github.com/kaplenko/diplom/internal/config"
	"github.com/kaplenko/diplom/internal/models"
)

func newTestJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		Secret:            "test-secret-key-for-testing",
		AccessExpiryHours: 1,
		RefreshExpiryDays: 7,
	}
}

func TestGenerateAndParseAccessToken(t *testing.T) {
	jwtCfg := newTestJWTConfig()
	svc := &AuthService{jwtCfg: jwtCfg}

	user := &models.User{
		ID:   42,
		Role: models.RoleStudent,
	}

	tokens, err := svc.generateTokenPair(user)
	if err != nil {
		t.Fatalf("generateTokenPair failed: %v", err)
	}

	if tokens.AccessToken == "" {
		t.Fatal("access token is empty")
	}
	if tokens.RefreshToken == "" {
		t.Fatal("refresh token is empty")
	}

	userID, role, err := svc.ParseAccessToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("ParseAccessToken failed: %v", err)
	}

	if userID != 42 {
		t.Errorf("expected user ID 42, got %d", userID)
	}
	if role != models.RoleStudent {
		t.Errorf("expected role student, got %s", role)
	}
}

func TestParseAccessToken_InvalidToken(t *testing.T) {
	jwtCfg := newTestJWTConfig()
	svc := &AuthService{jwtCfg: jwtCfg}

	_, _, err := svc.ParseAccessToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestParseAccessToken_RefreshTokenRejected(t *testing.T) {
	jwtCfg := newTestJWTConfig()
	svc := &AuthService{jwtCfg: jwtCfg}

	user := &models.User{ID: 1, Role: models.RoleAdmin}
	tokens, err := svc.generateTokenPair(user)
	if err != nil {
		t.Fatalf("generateTokenPair failed: %v", err)
	}

	_, _, err = svc.ParseAccessToken(tokens.RefreshToken)
	if err == nil {
		t.Fatal("expected error when using refresh token as access token")
	}
}

func TestGenerateTokenPair_AdminRole(t *testing.T) {
	jwtCfg := newTestJWTConfig()
	svc := &AuthService{jwtCfg: jwtCfg}

	user := &models.User{ID: 7, Role: models.RoleAdmin}

	tokens, err := svc.generateTokenPair(user)
	if err != nil {
		t.Fatalf("generateTokenPair failed: %v", err)
	}

	userID, role, err := svc.ParseAccessToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("ParseAccessToken failed: %v", err)
	}

	if userID != 7 {
		t.Errorf("expected user ID 7, got %d", userID)
	}
	if role != models.RoleAdmin {
		t.Errorf("expected role admin, got %s", role)
	}
}
