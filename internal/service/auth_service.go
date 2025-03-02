package service

import (
	"errors"
	"time"

	"blazperic/radionica/internal/models"
	"blazperic/radionica/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo                 *repository.UserRepository
	jwtSecret            string
	tokenDuration        time.Duration
	refreshTokenDuration time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string, tokenDuration, refreshTokenDuration time.Duration) *AuthService {
	return &AuthService{
		repo:                 repo,
		jwtSecret:            jwtSecret,
		tokenDuration:        tokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (s *AuthService) Register(username, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        uuid.New(),
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(username, password string) (*TokenPair, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.generateToken(user.ID, s.tokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, s.refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id not found in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id in token")
	}

	newAccessToken, err := s.generateToken(userID, s.tokenDuration)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateToken(userID, s.refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) generateToken(userID uuid.UUID, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(duration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
