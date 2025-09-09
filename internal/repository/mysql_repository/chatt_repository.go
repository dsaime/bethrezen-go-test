package mysqlRepository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"newsapi/internal/domain/newsAgr"
	sqlxRepo "newsapi/internal/repository/mysql_repository/sqlx_repo"
)

type NewsRepository struct {
	sqlxRepo.SqlxRepo
}

func (r *NewsRepository) Find(filter newsAgr.Filter) (newsAgr.News, error) {
	//TODO implement me
	panic("implement me")
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

	return toDomainNews(news, participantsMap, invitationsMap), nil
}

func (r *NewsRepository) Upsert(news newsAgr.News) (id int, _ error) {
	if news.ID < 1 {
		return 0, fmt.Errorf("chat ID is required")
	}

	if r.IsTx() {
		return r.upsert(news)
	} else {
		return id, r.InTransaction(func(txRepo newsAgr.Repository) error {
			var err error
			id, err = txRepo.Upsert(news)
			return err
		})
	}
}

func (r *NewsRepository) upsert(chat newsAgr.Chat) (int, error) {
	if _, err := r.DB().NamedExec(`
		INSERT INTO chats(id, name, chief_id) 
		VALUES (:id, :name, :chief_id)
		ON CONFLICT (id) DO UPDATE SET
			name=excluded.name,
			chief_id=excluded.chief_id
	`, toDBNews(chat)); err != nil {
		return fmt.Errorf("r.DB().NamedExec: %w", err)
	}

	// Удалить прошлых участников
	if _, err := r.DB().Exec(`
		DELETE FROM participants WHERE chat_id = $1
	`, chat.ID); err != nil {
		return fmt.Errorf("r.DB().Exec: %w", err)
	}

	if len(chat.Participants) > 0 {
		if _, err := r.DB().NamedExec(`
			INSERT INTO participants(chat_id, user_id)
			VALUES (:chat_id, :user_id)
		`, toDBParticipants(chat)); err != nil {
			return fmt.Errorf("r.DB().NamedExec: %w", err)
		}
	}

	// Удалить прошлые приглашения
	if _, err := r.DB().Exec(`
		DELETE FROM invitations WHERE chat_id = $1
	`, chat.ID); err != nil {
		return fmt.Errorf("r.DB().Exec: %w", err)
	}

	if len(chat.Invitations) > 0 {
		if _, err := r.DB().NamedExec(`
		INSERT INTO invitations(id, chat_id, subject_id, recipient_id)
		VALUES (:id, :chat_id, :subject_id, :recipient_id)
	`, toDBInvitations(chat)); err != nil {
			return fmt.Errorf("r.DB().NamedExec: %w", err)
		}
	}

	return nil
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
