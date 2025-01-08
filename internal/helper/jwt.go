package helper

import (
	"time"

	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type helper struct {
    config *viper.Viper
}

type AuthCustomClaims struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    jwt.StandardClaims
}

func NewHelper(config *viper.Viper) Helper {
    return &helper{
        config: config,
    }
}

func (h *helper) GenerateTokenUser(user model.UserResponse) (string, string, error) {
    refreshSecret := h.config.GetString("credentials.refreshsecret")
    accessSecret := h.config.GetString("credentials.accesssecret")

    accessTokenClaims := &AuthCustomClaims{
        Name:  user.Name,
        Email: user.Email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    refreshTokenClaims := &AuthCustomClaims{
        Name:  user.Name,
        Email: user.Email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
    accessTokenString, err := accessToken.SignedString([]byte(accessSecret))
    if err != nil {
        return "", "", err
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
    refreshTokenString, err := refreshToken.SignedString([]byte(refreshSecret))
    if err != nil {
        return "", "", err
    }

    return accessTokenString, refreshTokenString, nil
}