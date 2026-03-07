package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/abyalax/Boilerplate-go-gin/src/db"
	"github.com/abyalax/Boilerplate-go-gin/src/reject"
)

type UserService struct {
	q *Queries
}

func NewUserService(q *Queries) *UserService {
	return &UserService{q: q}
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (int32, error) {
	_, err := s.q.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return 0, reject.UserAlreadyExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	arg := CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	createdUser, err := s.q.CreateUser(ctx, arg)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return 0, reject.UserAlreadyExists
		}
		return 0, err
	}

	return int32(createdUser.ID), nil
}

func (s *UserService) GetUser(ctx context.Context, id int32) (*UserDTO, error) {

	u, err := s.q.GetUserByID(ctx, int32(id))
	if err != nil {
		if db.IsNoRows(err) {
			return nil, reject.UserNotFound
		}
		return nil, err
	}

	return &UserDTO{
		ID:    int32(u.ID),
		Name:  u.Name,
		Email: u.Email,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, req *UpdateUserRequest) (*UserDTO, error) {

	existing, err := s.q.GetUserByID(ctx, int32(id))
	if err != nil {
		if db.IsNoRows(err) {
			return nil, reject.UserNotFound
		}
		return nil, err
	}

	name := existing.Name
	email := existing.Email
	password := existing.Password

	if req.Name != nil && *req.Name != "" {
		name = *req.Name
	}

	if req.Email != nil && *req.Email != "" {
		email = *req.Email
		if email != existing.Email {
			_, err := s.q.GetUserByEmail(ctx, email)
			if err == nil {
				return nil, reject.UserAlreadyExists
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
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}

	updatedUser, err := s.q.UpdateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		ID:    int32(updatedUser.ID),
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {

	_, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		if db.IsNoRows(err) {
			return reject.UserNotFound
		}
		return err
	}

	err = s.q.DeleteUser(ctx, int32(id))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]UserDTO, error) {
	usersList, err := s.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	if len(usersList) == 0 {
		return []UserDTO{}, nil
	}

	dtos := make([]UserDTO, len(usersList))
	for i, u := range usersList {
		dtos[i] = UserDTO{
			ID:    int32(u.ID),
			Name:  u.Name,
			Email: u.Email,
		}
	}

	return dtos, nil
}
