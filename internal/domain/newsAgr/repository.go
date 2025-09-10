package newsAgr

import "errors"

var ErrNewsNotFound = errors.New("новость с такими параметрами не существует")

// Repository представляет собой интерфейс для работы с репозиторием новостей
type Repository interface {
	Find(Filter) (News, error)
	List(Filter) ([]News, error)
	Upsert(News) (id int, _ error)
	InTransaction(func(txRepo Repository) error) error
}

// Filter представляет собой фильтр для выборки новостей
type Filter struct {
	ID int // Фильтрация по ID
}
