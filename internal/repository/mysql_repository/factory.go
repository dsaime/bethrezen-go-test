package mysqlRepository

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"newsapi/internal/domain/newsAgr"
	sqlxRepo "newsapi/internal/repository/mysql_repository/sqlx_repo"
)

// Config представляет собой конфигурацию репозитория
type Config struct {
	DSN string
}

// Factory используется для создания репозиториев реализованных с помощью postgresql
type Factory struct {
	db *sqlx.DB
}

func InitFactory(cfg Config) (*Factory, error) {
	conn, err := sqlx.Connect("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("conn.Ping: %w", err)
	}

	logrus.Info("Успешно подключен к базе данных")

	return &Factory{
		db: conn,
	}, nil
}

// Close закрывает соединение с базой данных
func (f *Factory) Close() error {
	return f.db.Close()
}

// Cleanup очищает все сохраненные записи
func (f *Factory) Cleanup() error {
	var dbName string
	err := f.db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		return fmt.Errorf("не удалось получить имя базы данных: %w", err)
	}
	if !strings.HasPrefix(dbName, "test_") {
		return fmt.Errorf("очистка возможна только на тестовыз базах данных, текущая: %s", dbName)
	}

	// Список таблиц для очистки
	tables := []string{"News", "NewsCategories"}

	return sqlxRepo.New(f.db).InTransaction(func(tx sqlxRepo.SqlxRepo) error {
		for _, table := range tables {
			if _, err := tx.DB().Exec("DELETE FROM " + table); err != nil {
				return fmt.Errorf("tx.DB().Exec: %w", err)
			}
		}
		return nil
	})
}

// NewNewsRepository создает репозиторий чатов
func (f *Factory) NewNewsRepository() newsAgr.Repository {
	return &NewsRepository{
		SqlxRepo: sqlxRepo.New(f.db),
	}
}
