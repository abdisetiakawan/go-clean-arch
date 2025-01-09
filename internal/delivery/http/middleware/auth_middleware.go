package middleware

import (
	"encoding/json"
	"strings"

	"github.com/abdisetiakawan/go-clean-arch/internal/helper"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func NewAuth(userUserCase *usecase.UserUseCase, viper *viper.Viper, cache *helper.CacheHelper) fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        authHeader := ctx.Get("Authorization")
        if authHeader == "" {
            return fiber.ErrUnauthorized
        }
        secretkey := viper.GetString("credentials.accesssecret")
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secretkey), nil
        })
        if err != nil || !token.Valid {
            return fiber.ErrUnauthorized
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            return fiber.ErrUnauthorized
        }

        email := claims["email"].(string)
        sessionKey := "session:" + email
        sessionDataJSON, err := cache.Get(ctx.Context(), sessionKey)
        if err != nil {
            return fiber.ErrUnauthorized
        }

        var sessionData map[string]string
        if err := json.Unmarshal([]byte(sessionDataJSON), &sessionData); err != nil {
            return fiber.ErrUnauthorized
        }

        if sessionData["accessToken"] != tokenString {
            return fiber.ErrUnauthorized
        }

        auth := &model.Auth{
            Email: email,
        }
        ctx.Locals("auth", auth)
        return ctx.Next()
    }
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
    return ctx.Locals("auth").(*model.Auth)
}