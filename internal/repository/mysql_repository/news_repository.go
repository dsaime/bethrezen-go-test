package mysqlRepository

import (
	"fmt"
	"slices"

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
	// Запросить новости
	var news []dbNews
	if err := r.DB().Select(&news, `
		SELECT *
		FROM News
		WHERE (? = 0 OR Id = ?)
	`, filter.ID, filter.ID); err != nil {
		return nil, err
	}

	// Сразу вернуть пустой список
	if len(news) == 0 {
		return nil, nil
	}

	// Собрать ID новостей
	newsIDs := make([]any, len(news))
	for i := range news {
		newsIDs[i] = news[i].ID
	}

	// Создаем плейсхолдеры для IN clause
	placeholders := slices.Repeat([]rune("?,"), len(newsIDs))
	placeholders = placeholders[:len(placeholders)-1]

	// Запросить категории
	var categories []dbCategory
	if err := r.DB().Select(&categories, `
		SELECT *
		FROM NewsCategories
		WHERE NewsId IN (`+string(placeholders)+`)
	`, newsIDs...); err != nil {
		return nil, fmt.Errorf("r.DB().Select: %w", err)
	}

	categoriesMap := make(map[int][]dbCategory)
	for _, category := range categories {
		newsCategories := categoriesMap[category.NewsID]
		newsCategories = append(newsCategories, category)
		categoriesMap[category.NewsID] = newsCategories
	}

	return toDomainNews(news, categoriesMap), nil
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
		res, err := r.DB().Exec(`
			INSERT INTO News(Title, Content) 
			VALUES (?, ?)
		`, newsDb.Title, newsDb.Content)
		if err != nil {
			return 0, fmt.Errorf("r.DB().NamedExec: %w", err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("res.LastInsertId: %w", err)
		}
		newsDb.ID = int(id)
	} else {
		// Обновить запись
		if _, err := r.DB().NamedExec(`
			UPDATE News
			SET Title = :Title,
				Content = :Content
			WHERE Id = :Id
		`, newsDb); err != nil {
			return 0, fmt.Errorf("r.DB().NamedExec: %w", err)
		}
		// Удалить прошлые категории
		if _, err := r.DB().Exec(`
			DELETE FROM NewsCategories WHERE NewsId = ?
		`, news.ID); err != nil {
			return 0, fmt.Errorf("r.DB().Exec: %w", err)
		}
	}

	// Создать категории
	dbCategories := toDBCategories(newsDb.ID, news.Categories)
	if len(news.Categories) > 0 {
		if _, err := r.DB().NamedExec(`
			INSERT INTO NewsCategories(NewsId, CategoryId)
			VALUES (:NewsId, :CategoryId)
		`, dbCategories); err != nil {
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
	ID      int    `db:"Id"`
	Title   string `db:"Title"`
	Content string `db:"Content"`
	//Categories []int  `db:"Categories"`
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
			ID:     domainCategories[i],
			NewsID: newsID,
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

//func toNews2(
//	news dbNews,
//) newsAgr.News {
//	return newsAgr.News{
//		ID:         news.ID,
//		Title:      news.Title,
//		Content:    news.Content,
//		Categories: news.Categories,
//	}
//}
//
//func toDomainNews2(news []dbNews) []newsAgr.News {
//	domainNews := make([]newsAgr.News, len(news))
//	for i, news1 := range news {
//		domainNews[i] = toNews2(news1)
//	}
//
//	return domainNews
//}

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
