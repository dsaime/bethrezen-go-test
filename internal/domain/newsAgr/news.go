package newsAgr

import (
	"errors"
	"slices"
)

var (
	ErrEmptyTitle          = errors.New("пустой title")
	ErrEmptyContent        = errors.New("пустой Content")
	ErrDuplicateCategories = errors.New("повторяются категории")
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
	if err := ValidateTitle(title); err != nil {
		return News{}, err
	}
	if err := ValidateContent(content); err != nil {
		return News{}, err
	}
	if err := ValidateCategories(categories); err != nil {
		return News{}, err
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
		return ErrEmptyTitle
	}

	return nil
}

func ValidateContent(content string) error {
	if content == "" {
		return ErrEmptyContent
	}

	return nil
}

func ValidateCategories(categories []int) error {
	if len(categories) == 0 {
		return nil
	}

	if len(slices.Compact(categories)) != len(categories) {
		return ErrDuplicateCategories
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
