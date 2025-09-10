package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TokenVerifier описывает интерфейс проверки токена аутентификации
type TokenVerifier interface {
	VerifyToken(token string) bool
}

// RequireAuthorizedSession требует аутентификацию по Authorization заголовку
func RequireAuthorizedSession(verifier TokenVerifier) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Прочитать заголовок
		header := ctx.Get("Authorization")
		token, _ := strings.CutPrefix(header, "Bearer ")
		if token == "" {
			return fiber.ErrUnauthorized
		}

		// Проверить токен
		if !verifier.VerifyToken(token) {
			return fiber.ErrUnauthorized
		}

		return ctx.Next()
	}
}
