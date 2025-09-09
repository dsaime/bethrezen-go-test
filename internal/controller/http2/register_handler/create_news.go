package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	createNews "newsapi/internal/usecases/news/create_news"
)

// CreateNews регистрирует обработчик, позволяющий создать новость.
//
// Метод: POST /create
func CreateNews(router *fiber.App, uc UsecasesForCreateNews) {
	// Тело запроса для создания чата.
	type requestBody struct {
		Title      string `json:"Title"`
		Content    string `json:"Content"`
		Categories []int  `json:"Categories"`
	}
	router.Post(
		"/create",
		recover2.New(),
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			input := createNews.In{
				Title:      rb.Title,
				Content:    rb.Content,
				Categories: rb.Categories,
			}

			out, err := uc.CreateNews(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// UsecasesForCreateNews определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForCreateNews interface {
	CreateNews(createNews.In) (createNews.Out, error)
}
