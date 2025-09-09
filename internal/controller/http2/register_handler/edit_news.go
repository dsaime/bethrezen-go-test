package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	updateNews "newsapi/internal/usecases/news/update_news"
)

// EditNews регистрирует обработчик, позволяющий обновить параметры новости
//
// Метод: POST /edit/:Id
func EditNews(router *fiber.App, uc UsecasesForUpdateNews) {
	// Тело запроса для обновления названия чата.
	type requestBody struct {
		Title      string `json:"Title"`
		Content    string `json:"Content"`
		Categories []int  `json:"Categories"`
	}
	router.Post(
		"/edit/:Id",
		recover2.New(),
		func(ctx *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := ctx.BodyParser(&rb); err != nil {
				return err
			}

			input := updateNews.In{
				ID:         ParamsInt(ctx, "id"),
				Title:      rb.Title,
				Content:    rb.Content,
				Categories: rb.Categories,
			}

			out, err := uc.UpdateNews(input)
			if err != nil {
				return err
			}

			return ctx.JSON(out)
		},
	)
}

// UsecasesForUpdateNews определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForUpdateNews interface {
	UpdateNews(updateNews.In) (updateNews.Out, error)
}
