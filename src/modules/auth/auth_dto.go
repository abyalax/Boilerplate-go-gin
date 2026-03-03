package auth

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	User  GetUserWithPermissionsRow `json:"user"`
	Token string                    `json:"token"`
}
