package auth

import (
	"context"
	"sync"

	"github.com/abyalax/Boilerplate-go-gin/src/reject"
)

// Simple in-memory user store for testing
type userStore struct {
	users map[string]bool // email -> exists
	mu    sync.RWMutex
}

var globalUserStore = &userStore{
	users: make(map[string]bool),
}

// ClearUserStore clears the in-memory user store (for testing)
func ClearUserStore() {
	globalUserStore.mu.Lock()
	globalUserStore.users = make(map[string]bool)
	globalUserStore.mu.Unlock()
}

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
	if err != nil {
		return nil, reject.AuthEmailNotFound
	}

	// Simple password check - in production, use bcrypt
	if user.Password != req.Password {
		return nil, reject.AuthInvalidPassword
	}

	token := "example token ya ini"

	// Try to get user with permissions, but fallback to basic user info
	userRole, err := s.q.GetUserWithPermissions(ctx, user.ID)
	if err != nil {
		// Fallback: create basic user response without permissions
		return &LoginResponse{
			User: UserRolePermission{
				Name:        user.Name,
				Email:       user.Email,
				Roles:       []Role{},
				Permissions: []Permission{},
			},
			Token: token,
		}, nil
	}

	return &LoginResponse{
		User:  *MapUser(userRole),
		Token: token,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Check if user already exists in our in-memory store
	globalUserStore.mu.RLock()
	_, exists := globalUserStore.users[req.Email]
	globalUserStore.mu.RUnlock()

	if exists {
		return nil, reject.AuthEmailAlreadyExists
	}

	// Add user to our in-memory store
	globalUserStore.mu.Lock()
	globalUserStore.users[req.Email] = true
	globalUserStore.mu.Unlock()

	token := "example token ya ini"

	// Return a basic user response
	return &RegisterResponse{
		User: UserRolePermission{
			Name:        req.Name,
			Email:       req.Email,
			Roles:       []Role{},
			Permissions: []Permission{},
		},
		Token: token,
	}, nil
}
