package createNews

import (
	"newsapi/internal/domain/newsAgr"
)

// In входящие параметры
type In struct {
	Title      string
	Content    string
	Categories []int
}

// Out результат создания агрегата
type Out struct {
	News newsAgr.News
}

type CreateNewsUsecase struct {
	Repo newsAgr.Repository
}

// CreateNews создает новость
func (c *CreateNewsUsecase) CreateNews(in In) (Out, error) {
	news, err := newsAgr.NewNews(c.Repo, in.Title, in.Content, in.Categories)
	if err != nil {
		return Out{}, err
	}

	return Out{
		News: news,
	}, nil
}
