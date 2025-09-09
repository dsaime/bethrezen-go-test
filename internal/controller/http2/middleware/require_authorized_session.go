package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthTokenVerifier описывает интерфейс проверки токена аутентификации
type AuthTokenVerifier interface {
	VerifyAuthToken(token string) bool
}

// RequireAuthorizedSession требует аутентификацию по Authorization заголовку
func RequireAuthorizedSession(verifier AuthTokenVerifier) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Прочитать заголовок
		header := ctx.Get("Authorization")
		token, _ := strings.CutPrefix(header, "Bearer ")
		if token == "" {
			return fiber.ErrUnauthorized
		}

		// Проверить токен
		if verifier.VerifyAuthToken(token) {
			return fiber.ErrUnauthorized
		}

		return ctx.Next()
	}
}
