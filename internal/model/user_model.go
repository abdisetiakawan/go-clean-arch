package model

type UserResponse struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}

type CreateUserRequest struct {
	Name     string `json:"name,omitempty" validate:"required,max=100"`
	Email    string `json:"email,omitempty" validate:"required,max=100"`
	Password string `json:"password,omitempty" validate:"required,max=100"`
}

type LoginUserRequest struct {
	Email    string `json:"email,omitempty" validate:"required,max=100"`
	Password string `json:"password,omitempty" validate:"required,max=100"`
}

type LogoutUserRequest struct {
	Email string `json:"email,omitempty" validate:"required,max=100"`
}

type UpdateUserRequest struct {
	Name     string `json:"name,omitempty" validate:"required,max=100"`
	Email    string `json:"email,omitempty" validate:"required,max=100"`
	Password string `json:"password,omitempty" validate:"required,max=100"`
}

type GetUserRequest struct {
	Email string `json:"email,omitempty" validate:"required,max=100"`
}