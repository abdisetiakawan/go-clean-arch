package middleware

import (
	"strings"

	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func NewAuth(userUserCase *usecase.UserUseCase, viper *viper.Viper) fiber.Handler {
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

        auth := &model.Auth{
            Email: claims["email"].(string),
        }
        ctx.Locals("auth", auth)
        return ctx.Next()
    }
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
