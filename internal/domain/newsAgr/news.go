package newsAgr

import (
	"errors"
	"slices"
	"unicode"
)

var (
	ErrInvalidTitle      = errors.New("некорректное значение Title")
	ErrInvalidContent    = errors.New("некорректное значение Content")
	ErrInvalidCategories = errors.New("некорректное значение Categories")
)

// News представляет собой агрегат новости
type News struct {
	ID         int    // Уникальный ID новости
	Title      string // Название
	Content    string // Содержимое новости
	Categories []int  // Список категорий
}

// NewNews создает новую новость
func NewNews(repo Repository, title string, content string, categories []int) (News, error) {
	errs := errors.Join(
		ValidateTitle(title),
		ValidateContent(content),
		ValidateCategories(categories),
	)
	if errs != nil {
		return News{}, errs
	}

	news := News{
		ID:         0,
		Title:      title,
		Content:    content,
		Categories: categories,
	}

	id, err := repo.Upsert(news)
	if err != nil {
		return News{}, err
	}
	news.ID = id

	return news, nil
}

func ValidateTitle(title string) error {
	if title == "" {
		return ValidationError{
			err:    ErrInvalidTitle,
			reason: "пустой заголовок",
		}
	}

	firstRune := []rune(title)[0]
	if unicode.Is(unicode.White_Space, firstRune) {
		return ValidationError{
			err:    ErrInvalidTitle,
			reason: "начинаться с пробела",
		}
	}
	if unicode.Is(unicode.Lower, firstRune) {
		return ValidationError{
			err:    ErrInvalidTitle,
			reason: "начинается с символа нижнего регистра",
		}
	}

	return nil
}

func ValidateContent(content string) error {
	if content == "" {
		return ValidationError{
			err:    ErrInvalidContent,
			reason: "пустой",
		}
	}

	return nil
}

func ValidateCategories(categories []int) error {
	if len(categories) == 0 {
		return nil
	}

	if len(slices.Compact(categories)) != len(categories) {
		return ValidationError{
			err:    ErrInvalidCategories,
			reason: "повторяющиеся категории",
		}
	}

	return nil
}

func (n News) UpdateTitle(repo Repository, title string) (News, error) {
	if err := ValidateTitle(title); err != nil {
		return News{}, err
	}

	newNews := News{
		ID:         n.ID,
		Title:      title,
		Content:    n.Content,
		Categories: n.Categories,
	}
	if _, err := repo.Upsert(newNews); err != nil {
		return News{}, err
	}

	return newNews, nil
}

func (n News) UpdateContent(repo Repository, content string) (News, error) {
	if err := ValidateContent(content); err != nil {
		return News{}, err
	}

	newNews := News{
		ID:         n.ID,
		Title:      n.Title,
		Content:    content,
		Categories: n.Categories,
	}
	if _, err := repo.Upsert(newNews); err != nil {
		return News{}, err
	}

	return newNews, nil
}

func (n News) UpdateCategories(repo Repository, categories []int) (News, error) {
	if err := ValidateCategories(categories); err != nil {
		return News{}, err
	}

	newNews := News{
		ID:         n.ID,
		Title:      n.Title,
		Content:    n.Content,
		Categories: categories,
	}
	if _, err := repo.Upsert(n); err != nil {
		return News{}, err
	}

	return newNews, nil
}
