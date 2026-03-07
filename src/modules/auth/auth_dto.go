package auth

import "database/sql"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	User  UserRolePermission `json:"user"`
	Token string             `json:"token"`
}

type RegisterResponse struct {
	User  UserRolePermission `json:"user"`
	Token string             `json:"token"`
}

type GetUserAuthRow struct {
	UserID         int64
	UserName       string
	UserEmail      string
	RoleID         sql.NullInt64
	RoleName       sql.NullString
	PermissionID   sql.NullInt64
	PermissionKey  sql.NullString
	PermissionName sql.NullString
}
