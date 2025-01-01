package model

type UserResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Token     string `json:"token,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}

type CreateUserRequest struct {
	ID       string `json:"id,omitempty" validate:"required,max=100"`
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