package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"newsapi/internal/controller/http2/middleware"
	newsList "newsapi/internal/usecases/news/news_list"
)

// NewsList регистрирует HTTP-обработчик для получения списка новостей
//
// Метод: GET /list
func NewsList(router *fiber.App, uc UsecasesForNewsList, verifier middleware.TokenVerifier) {
	router.Get(
		"/list",
		recover2.New(),
		middleware.RequireAuthorizedSession(verifier),
		func(ctx *fiber.Ctx) error {
			out, err := uc.NewsList(newsList.In{})
			if err != nil {
				return err
			}

			return ctx.JSON(fiber.Map{
				"Success": true,
				"News":    out.News,
			})
		},
	)
}

// UsecasesForNewsList определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForNewsList interface {
	NewsList(newsList.In) (newsList.Out, error)
}
