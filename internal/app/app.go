package app

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"newsapi/internal/controller/http2"
)

func Run(ctx context.Context, cfg Config) error {
	// Установить уровень логирования
	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("logrus.ParseLevel: %w", err)
	}
	logrus.SetLevel(logLevel)
	logrus.Info(fmt.Sprintf("Уровень логирования: %s", cfg.LogLevel))

	// Инициализация репозиториев
	rr, closeRepos, err := initPgsqlRepositories(cfg.Pgsql)
	if err != nil {
		return err
	}
	defer closeRepos()

	// Инициализация адаптеров
	aa := initAdapters(cfg)

	// Инициализация сервисов
	uc := initUsecases(rr, aa)

	return http2.RunHttpServer(ctx, uc, aa.eventBus, cfg.Http2)
}
