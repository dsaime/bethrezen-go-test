package app

import (
	"fmt"

	"newsapi/internal/domain/newsAgr"
	mysqlRepository "newsapi/internal/repository/mysql_repository"
)

// repositories содержит инициализированный репозитории
type repositories struct {
	factory  *mysqlRepository.Factory
	newsRepo newsAgr.Repository
}

func (r *repositories) Close() error {
	if err := r.factory.Close(); err != nil {
		return fmt.Errorf("r.factory.Close: %w", err)
	}

	return nil
}

// initMysqlRepositories инициализирует репозитории, реализованные с помощью mysql
func initMysqlRepositories(cfg mysqlRepository.Config) (*repositories, error) {
	// Создать фабрику репозиториев
	factory, err := mysqlRepository.InitFactory(cfg)
	if err != nil {
		return nil, fmt.Errorf("mysqlRepository.InitFactory: %w", err)
	}

	rs := &repositories{
		newsRepo: factory.NewNewsRepository(),
		factory:  factory,
	}

	return rs, nil
}
