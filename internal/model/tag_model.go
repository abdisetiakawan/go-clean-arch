package model

type CreateTagRequest struct {
	Email string `json:"email" validate:"required,max=150"`
	Name  string `json:"name" validate:"required,max=50"`
}

type TagResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name"`
}

type SearchTagRequest struct {
	Email string `json:"-"`
	Name  string `json:"name"`
	Page  int    `json:"page" validate:"min=1"`
	Size  int    `json:"size" validate:"min=1,max=100"`
}

type GetTagRequest struct {
	ID    string `json:"-" validate:"required"`
	Email string `json:"-" validate:"required"`
}

type UpdateTagRequest struct {
	ID			string `json:"-"`
	Email       string `json:"-" validate:"max=100"`
	Name       string `json:"name" validate:"max=150"`
}