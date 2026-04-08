package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kaplenko/diplom/internal/config"
	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtCfg   config.JWTConfig
}

func NewAuthService(userRepo *repository.UserRepository, jwtCfg config.JWTConfig) *AuthService {
	return &AuthService{userRepo: userRepo, jwtCfg: jwtCfg}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.User, error) {
	existing, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if existing != nil {
		return nil, ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		Role:         models.RoleStudent,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.TokenResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrUnauthorized
	}

	return s.generateTokenPair(user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.TokenResponse, error) {
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		return nil, ErrInvalidToken
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(int64(userIDFloat))
	if err != nil || user == nil {
		return nil, ErrInvalidToken
	}

	return s.generateTokenPair(user)
}

// ParseAccessToken validates an access token and returns the user ID and role.
func (s *AuthService) ParseAccessToken(tokenStr string) (int64, models.Role, error) {
	claims, err := s.parseToken(tokenStr)
	if err != nil {
		return 0, "", ErrInvalidToken
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "access" {
		return 0, "", ErrInvalidToken
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", ErrInvalidToken
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", ErrInvalidToken
	}

	return int64(userIDFloat), models.Role(role), nil
}

func (s *AuthService) generateTokenPair(user *models.User) (*models.TokenResponse, error) {
	accessExpiry := time.Now().Add(time.Duration(s.jwtCfg.AccessExpiryHours) * time.Hour)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    string(user.Role),
		"type":    "access",
		"exp":     accessExpiry.Unix(),
		"iat":     time.Now().Unix(),
	})

	accessStr, err := accessToken.SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshExpiry := time.Now().Add(time.Duration(s.jwtCfg.RefreshExpiryDays) * 24 * time.Hour)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     refreshExpiry.Unix(),
		"iat":     time.Now().Unix(),
	})

	refreshStr, err := refreshToken.SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		ExpiresAt:    accessExpiry.Format(time.RFC3339),
	}, nil
}

func (s *AuthService) parseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtCfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
