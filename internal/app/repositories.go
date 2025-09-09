package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"newsapi/internal/domain/newsAgr"
	mysqlRepository "newsapi/internal/repository/mysql_repository"
)

type repositories struct {
	newsRepo newsAgr.Repository
}

func initPgsqlRepositories(cfg mysqlRepository.Config) (*repositories, func(), error) {
	factory, err := mysqlRepository.InitFactory(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("pgsqlRepository.InitFactory: %w", err)
	}

	rs := &repositories{
		newsRepo: factory.NewNewsRepository(),
	}

	closer := func() {
		if err := factory.Close(); err != nil {
			logrus.Error("Закрыть соединение с pgsql: factory.Close: " + err.Error())
		}
	}

	return rs, closer, nil
}
