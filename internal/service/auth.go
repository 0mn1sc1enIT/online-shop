package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/OmniscienIT/GolangAPI/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      repository.Users
	tokenTTL  time.Duration
	signedKey string
}

func NewAuthService(repo repository.Users, signedKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:      repo,
		signedKey: signedKey,
		tokenTTL:  tokenTTL,
	}
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

func (s *AuthService) CreateUser(user domain.User) (uint, error) {
	// 1. Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(passwordHash)

	// 2. Сохраняем в БД
	if err := s.repo.Create(&user); err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (s *AuthService) GenerateToken(email, password string) (string, error) {
	// 1. Ищем пользователя
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// 2. Сверяем хеши паролей
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// 3. Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
		Role:   user.Role,
	})

	return token.SignedString([]byte(s.signedKey))
}

func (s *AuthService) ParseToken(accessToken string) (uint, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(s.signedKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", fmt.Errorf("token claims are not of type *tokenClaims")
	}

	return claims.UserID, claims.Role, nil
}
