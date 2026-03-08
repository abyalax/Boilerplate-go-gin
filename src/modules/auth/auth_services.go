package auth

import (
	"context"
	"time"

	"github.com/abyalax/Boilerplate-go-gin/src/config/app"
	"github.com/abyalax/Boilerplate-go-gin/src/config/env"
	"github.com/abyalax/Boilerplate-go-gin/src/http"
	"github.com/abyalax/Boilerplate-go-gin/src/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	q *Queries
}

func NewAuthService(q *Queries) *AuthService {
	return &AuthService{
		q: q,
	}
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.q.GetUserByEmail(ctx, req.Email)
	cfg, _ := env.Load()
	if err != nil {
		return nil, app.Reject(http.AuthEmailNotFound, nil)
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, app.Reject(http.AuthInvalidPassword, nil)
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"iat":   now.Unix(),
		"exp":   now.Add(cfg.JWT.TokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return nil, app.Reject(http.JWTFailedGenerateToken, err)
	}

	userRole, err := s.q.GetUserWithPermissions(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:  *MapUser(userRole),
		Token: tokenString,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {

	_, err := s.q.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, app.Reject(http.AuthEmailAlreadyExists, nil)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.q.CreateUser(ctx, CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})

	return &RegisterResponse{
		User: User{
			ID:    user.ID,
			Name:  req.Name,
			Email: req.Email,
		},
	}, nil
}
