package services

import (
	"context"
	"errors"
	"regexp"
	"time"

	"cashier/models"
	"cashier/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) error
	Login(ctx context.Context, req models.AuthRequest) (string, error)
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
}

type authService struct {
	repo      repositories.AuthRepository
	jwtSecret []byte
}

func NewAuthService(repo repositories.AuthRepository, secret string) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: []byte(secret),
	}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) error {
	// 1. Validation
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Status:   "member", // Default status
	}

	return s.repo.Register(ctx, user)
}

func (s *authService) Login(ctx context.Context, req models.AuthRequest) (string, error) {
	user, err := s.repo.FindByUsernameOrEmail(ctx, req.Identity)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"status":  user.Status,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}

func (s *authService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}
