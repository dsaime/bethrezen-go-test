package middleware

import (
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequireAuthorizedSession(t *testing.T) {
	// Настройка мока
	trustedTokens := []string{"a", "b"}
	verifier := NewMockTokenVerifier(t)
	verifier.On("VerifyToken", mock.Anything).
		Maybe().
		Return(func(token string) bool {
			return slices.Contains(trustedTokens, token)
		})
	// Настройка сервера
	fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
	fiberApp.Get("/", RequireAuthorizedSession(verifier), func(ctx *fiber.Ctx) error {
		return ctx.SendString("ok")
	})

	t.Run("Без токена вернет ошибку авторизации", func(t *testing.T) {
		// Создать запрос
		req := httptest.NewRequest("GET", "/", nil)

		// Выполнить запрос
		response, err := fiberApp.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
	})

	t.Run("Невалидный токен вернет ошибку авторизации", func(t *testing.T) {
		// Создать запрос
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid")

		// Выполнить запрос
		response, err := fiberApp.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
	})

	t.Run("Валидный токен пройдет авторизацию", func(t *testing.T) {
		// Создать запрос
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", trustedTokens[0])

		// Выполнить запрос
		response, err := fiberApp.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)
	})
}
