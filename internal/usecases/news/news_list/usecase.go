package newsList

import (
	"newsapi/internal/domain/newsAgr"
)

// In входящие параметры
type In struct{}

// Out результат запроса новостей
type Out struct {
	News []newsAgr.News
}

type NewsListUsecase struct {
	Repo newsAgr.Repository
}

// NewsList возвращает список новостей
func (c *NewsListUsecase) NewsList(_ In) (Out, error) {
	// Получить список новостей
	news, err := c.Repo.List(newsAgr.Filter{})
	if err != nil {
		return Out{}, err
	}

	return Out{News: news}, err
}
