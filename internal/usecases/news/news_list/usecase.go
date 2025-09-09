package newsList

import (
	"newsapi/internal/domain/newsAgr"
)

// In входящие параметры
type In struct{}

// Out результат запроса чатов
type Out struct {
	News []newsAgr.News
}

type NewsListUsecase struct {
	Repo newsAgr.Repository
}

// NewsList возвращает список чатов, в которых участвует пользователь
func (c *NewsListUsecase) NewsList(_ In) (Out, error) {
	// Получить список новостей
	chats, err := c.Repo.List(newsAgr.Filter{})
	if err != nil {
		return Out{}, err
	}

	return Out{News: chats}, err
}
