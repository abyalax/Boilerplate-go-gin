package auth

import (
	"context"
)

type AuthService struct {
	q *Queries
}

func NewAuthService(q *Queries) *AuthService {
	return &AuthService{q: q}
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.q.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrEmailNotFound
	}

	token := "example token ya ini"

	userRole, err := s.q.GetUserWithPermissions(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:  userRole,
		Token: token,
	}, nil
}
