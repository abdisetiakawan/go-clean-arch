package helper

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
)

type Helper interface {
    GenerateTokenUser(user model.UserResponse) (string, string, error)
}