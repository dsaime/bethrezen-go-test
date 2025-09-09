package updateNews

import (
	"errors"

	"newsapi/internal/domain/newsAgr"
)

var (
	ErrInvalidNewsID   = errors.New("некорректное значение ID новости")
	ErrNothingToUpdate = errors.New("нечего обновлять")
)

// In входящие параметры
type In struct {
	ID         int
	Title      string
	Content    string
	Categories []int
}

// Validate валидирует значение отдельно каждого параметры
func (in In) Validate() error {
	if in.ID < 1 {
		return ErrInvalidNewsID
	}

	if in.isAllFieldEmpty() {
		return ErrNothingToUpdate
	}

	if in.Title != "" {
		if err := newsAgr.ValidateTitle(in.Title); err != nil {
			return err
		}
	}

	if in.Content != "" {
		if err := newsAgr.ValidateContent(in.Content); err != nil {
			return err
		}
	}

	if len(in.Categories) != 0 {
		if err := newsAgr.ValidateCategories(in.Categories); err != nil {
			return err
		}
	}

	return nil
}

func (in In) isAllFieldEmpty() bool {
	return in.Title == "" &&
		in.Content == "" &&
		len(in.Categories) == 0
}

// Out результат обновления названия чата
type Out struct {
	News newsAgr.News
}

type UpdateNewsUsecase struct {
	Repo newsAgr.Repository
}

// UpdateNews обновляет новость
func (c *UpdateNewsUsecase) UpdateNews(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Найти новость
	oldNews, err := c.Repo.Find(newsAgr.Filter{ID: in.ID})
	if err != nil {
		return Out{}, err
	}

	// Обновляемые новости
	newNews := oldNews
	// Выполнить обновление в транзакции
	if err = c.Repo.InTransaction(func(txRepo newsAgr.Repository) error {
		if in.Title != "" {
			if newNews, err = newNews.UpdateTitle(c.Repo, in.Title); err != nil {
				return err
			}
		}
		if in.Content != "" {
			if newNews, err = newNews.UpdateContent(c.Repo, in.Content); err != nil {
				return err
			}
		}
		if len(in.Categories) != 0 {
			if newNews, err = newNews.UpdateCategories(c.Repo, in.Categories); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return Out{}, err
	}

	return Out{News: newNews}, nil
}
