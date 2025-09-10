package app

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"newsapi/internal/controller/http2"
)

// Run запускает приложение
func Run(ctx context.Context, cfg Config) error {
	// Установить уровень логирования
	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("logrus.ParseLevel: %w", err)
	}
	logrus.SetLevel(logLevel)
	logrus.Info(fmt.Sprintf("Уровень логирования: %s", cfg.LogLevel))

	// Инициализация репозиториев
	rr, err := initMysqlRepositories(cfg.Mysql)
	if err != nil {
		return err
	}
	defer func() {
		if err := rr.Close(); err != nil {
			logrus.Error("Закрыть репозитории: ", err)
		}
	}()

	// Инициализация адаптеров
	aa := initAdapters(cfg)

	// Инициализация сервисов
	uc := initUsecases(rr)

	// Запуск http сервера
	return http2.RunHttpServer(ctx, cfg.Http2, uc, aa.TokenVerifier)
}
