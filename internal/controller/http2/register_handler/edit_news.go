package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"

	"newsapi/internal/controller/http2/middleware"
	updateNews "newsapi/internal/usecases/news/update_news"
)

// EditNews регистрирует обработчик, позволяющий обновить параметры новости
//
// Метод: POST /edit/:Id
func EditNews(router *fiber.App, uc UsecasesForUpdateNews, verifier middleware.TokenVerifier) {
	// Тело запроса
	type requestBody struct {
		Title      string `json:"Title"`
		Content    string `json:"Content"`
		Categories []int  `json:"Categories"`
	}
	router.Post(
		"/edit/:Id",
		recover2.New(),
		middleware.RequireAuthorizedSession(verifier),
		func(ctx *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := ctx.BodyParser(&rb); err != nil {
				return err
			}

			logOldNewsFields := logrus.Fields{
				"id":         ParamsInt(ctx, "id"),
				"title":      rb.Title,
				"content":    rb.Content,
				"categories": rb.Categories,
			}

			input := updateNews.In{
				ID:         ParamsInt(ctx, "id"),
				Title:      rb.Title,
				Content:    rb.Content,
				Categories: rb.Categories,
			}

			out, err := uc.UpdateNews(input)
			if err != nil {
				logrus.
					WithFields(logOldNewsFields).
					Warn("Ошибка при обновлении новости: ", err)
				return err
			}

			logrus.
				WithFields(logOldNewsFields).
				WithField("new_news", out.News).
				Info("Успешное обновление новости")

			return ctx.JSON(fiber.Map{
				"Success": true,
				"News":    out.News,
			})
		},
	)
}

// UsecasesForUpdateNews определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForUpdateNews interface {
	UpdateNews(updateNews.In) (updateNews.Out, error)
}
