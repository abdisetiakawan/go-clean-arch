package converter

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/model"

	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
)

func UserToResponse(user *entity.User) *model.UserResponse {
    return &model.UserResponse{
        ID:           user.ID,
        Name:         user.Name,
        Email:        user.Email,
        AccessToken:  user.AccessToken,
        RefreshToken: user.Token,
    }
}


func UserToTokenResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		AccessToken: user.Token,
	}
}
