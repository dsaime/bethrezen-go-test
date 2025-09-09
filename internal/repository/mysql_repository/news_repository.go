package mysqlRepository

import (
	"fmt"

	"newsapi/internal/domain/newsAgr"
	sqlxRepo "newsapi/internal/repository/mysql_repository/sqlx_repo"
)

type NewsRepository struct {
	sqlxRepo.SqlxRepo
}

func (r *NewsRepository) Find(filter newsAgr.Filter) (newsAgr.News, error) {
	news, err := r.List(filter)
	if err != nil {
		return newsAgr.News{}, err
	}
	if len(news) != 1 {
		return newsAgr.News{}, newsAgr.ErrNewsNotFound
	}

	return news[0], nil
}

func (r *NewsRepository) List(filter newsAgr.Filter) ([]newsAgr.News, error) {
	// Запросить чаты
	var news []dbNews
	if err := r.DB().Select(&news, `
		SELECT n.*, GROUP_CONCAT(n.title) as Categories
		FROM News n
			LEFT JOIN NewsCategories c
				ON n.Id = c.NewsId
		WHERE ($1 = 0 OR c.NewsId = $1)
	`, filter.ID); err != nil {
		return nil, err
	}

	// Сразу вернуть пустой список
	if len(news) == 0 {
		return nil, nil
	}

	return toDomainNews2(news), nil
}

func (r *NewsRepository) Upsert(news newsAgr.News) (int, error) {
	if r.IsTx() {
		return r.upsert(news)
	} else {
		var id int
		return id, r.InTransaction(func(txRepo newsAgr.Repository) error {
			var err error
			id, err = txRepo.Upsert(news)
			return err
		})
	}
}

func (r *NewsRepository) upsert(news newsAgr.News) (int, error) {
	newsDb := toDBNews(news)

	if newsDb.ID == 0 {
		// Вставить запись и получить её ID.
		// В newsDb.ID записывает ID созданной новости
		if err := r.DB().Get(&newsDb.ID, `
			INSERT INTO News(Title, Content) 
			VALUES ($1, $2)
			RETURNING Id
		`, newsDb.Title, newsDb.Content); err != nil {
			return 0, fmt.Errorf("r.DB().NamedExec: %w", err)
		}
	} else {
		// Обновить запись
		if _, err := r.DB().Exec(`
			UPDATE NewsCategories
			SET Title = :Title,
				Content = :Content	
			WHERE NewsId = :NewsId
		`); err != nil {
			return 0, fmt.Errorf("r.DB().NamedExec: %w", err)
		}
		// Удалить прошлые категории
		if _, err := r.DB().Exec(`
			DELETE FROM NewsCategories WHERE NewsId = $1
		`, news.ID); err != nil {
			return 0, fmt.Errorf("r.DB().Exec: %w", err)
		}
	}

	// Создать категории
	if len(news.Categories) > 0 {
		if _, err := r.DB().NamedExec(`
			INSERT INTO NewsCategories(NewsId, CategoryId)
			VALUES (:chat_id, :user_id)
		`, toDBCategories(news.ID, newsDb.Categories)); err != nil {
			return 0, fmt.Errorf("r.DB().NamedExec: %w", err)
		}
	}

	return newsDb.ID, nil
}

func (r *NewsRepository) InTransaction(fn func(txRepo newsAgr.Repository) error) error {
	return r.SqlxRepo.InTransaction(func(txSqlxRepo sqlxRepo.SqlxRepo) error {
		return fn(&NewsRepository{SqlxRepo: txSqlxRepo})
	})
}

type dbNews struct {
	ID         int    `db:"Id"`
	Title      string `db:"Title"`
	Content    string `db:"Content"`
	Categories []int  `db:"Categories"`
}

func toDBNews(news newsAgr.News) dbNews {
	return dbNews{
		ID:      news.ID,
		Title:   news.Title,
		Content: news.Content,
	}
}

func toDBCategories(newsID int, domainCategories []int) []dbCategory {
	categories := make([]dbCategory, len(domainCategories))
	for i := range domainCategories {
		categories[i] = dbCategory{
			ID:     newsID,
			NewsID: domainCategories[i],
		}
	}

	return categories
}

func toNews(
	news dbNews,
	categories []dbCategory,
) newsAgr.News {
	categoryIDs := make([]int, len(categories))
	for i := range categories {
		categoryIDs[i] = categories[i].ID
	}

	return newsAgr.News{
		ID:         news.ID,
		Title:      news.Title,
		Content:    news.Content,
		Categories: categoryIDs,
	}
}

func toNews2(
	news dbNews,
) newsAgr.News {
	return newsAgr.News{
		ID:         news.ID,
		Title:      news.Title,
		Content:    news.Content,
		Categories: news.Categories,
	}
}

func toDomainNews2(news []dbNews) []newsAgr.News {
	domainNews := make([]newsAgr.News, len(news))
	for i, news1 := range news {
		domainNews[i] = toNews2(news1)
	}

	return domainNews
}

func toDomainNews(
	news []dbNews,
	categoriesMap map[int][]dbCategory,
) []newsAgr.News {
	domainNews := make([]newsAgr.News, len(news))
	for i, news1 := range news {
		domainNews[i] = toNews(news1, categoriesMap[news1.ID])
	}

	return domainNews
}

type dbCategory struct {
	ID     int `db:"CategoryId"`
	NewsID int `db:"NewsId"`
}
