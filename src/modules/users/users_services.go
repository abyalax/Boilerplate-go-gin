package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

// UserService handles all user business logic
type UserService struct {
	q *Queries
}

// NewUserService creates a new UserService
func NewUserService(q *Queries) *UserService {
	return &UserService{q: q}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (int64, error) {

	// Check if user with this email already exists
	_, err := s.q.GetUserByEmail(ctx, req.Email)
	if err == nil {
		// User exists
		return 0, ErrUserAlreadyExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		// Other database error
		return 0, err
	}

	// Insert new user
	arg := CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	createdUser, err := s.q.CreateUser(ctx, arg)
	if err != nil {
		return 0, err
	}

	return int64(createdUser.ID), nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, idStr string) (*UserDTO, error) {
	id, err := s.parseUserID(idStr)
	if err != nil {
		return nil, err
	}

	u, err := s.q.GetUserByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &UserDTO{
		ID:    int64(u.ID),
		Name:  u.Name,
		Email: u.Email,
	}, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, idStr string, req *UpdateUserRequest) (*UserDTO, error) {
	id, err := s.parseUserID(idStr)
	if err != nil {
		return nil, err
	}

	// Fetch existing user
	existing, err := s.q.GetUserByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	name := existing.Name
	email := existing.Email
	password := existing.Password

	// Update fields if provided
	if req.Name != nil && *req.Name != "" {
		name = *req.Name
	}
	if req.Email != nil && *req.Email != "" {
		email = *req.Email
		// Check for duplicate email only if email is being changed
		if email != existing.Email {
			_, err := s.q.GetUserByEmail(ctx, email)
			if err == nil {
				// User with this email already exists
				return nil, ErrUserAlreadyExists
			}
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
	}
	if req.Password != nil && *req.Password != "" {
		password = *req.Password
	}

	arg := UpdateUserParams{
		ID:       int32(id),
		Name:     name,
		Email:    email,
		Password: password,
	}
	updatedUser, err := s.q.UpdateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		ID:    int64(updatedUser.ID),
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, idStr string) error {
	id, err := s.parseUserID(idStr)
	if err != nil {
		return err
	}

	// Check if user exists first
	_, err = s.q.GetUserByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	err = s.q.DeleteUser(ctx, int32(id))
	if err != nil {
		return err
	}

	return nil
}

// ListUsers retrieves all users
func (s *UserService) ListUsers(ctx context.Context) ([]UserDTO, error) {
	usersList, err := s.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	// Return empty slice instead of nil
	if len(usersList) == 0 {
		return []UserDTO{}, nil
	}

	dtos := make([]UserDTO, len(usersList))
	for i, u := range usersList {
		dtos[i] = UserDTO{
			ID:    int64(u.ID),
			Name:  u.Name,
			Email: u.Email,
		}
	}

	return dtos, nil
}

// UserDTO is the standard user response DTO
type UserDTO struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// parseUserID parses and validates user ID from path parameter
func (s *UserService) parseUserID(idStr string) (int64, error) {
	if idStr == "" {
		return 0, fmt.Errorf("user id is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user id: must be a valid number")
	}

	if id <= 0 {
		return 0, fmt.Errorf("invalid user id: must be greater than 0")
	}

	return id, nil
}
