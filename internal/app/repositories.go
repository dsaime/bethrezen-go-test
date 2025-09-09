package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"newsapi/internal/domain/chatt"
	"newsapi/internal/domain/sessionn"
	"newsapi/internal/domain/userr"
	pgsqlRepository "newsapi/internal/repository/pgsql_repository"

	"newsapi/internal/domain/newsAgr"
)

type repositories struct {
	chats    newsAgr.Repository
	users    userr.Repository
	sessions sessionn.Repository
}

func initPgsqlRepositories(cfg pgsqlRepository.Config) (*repositories, func(), error) {
	factory, err := pgsqlRepository.InitFactory(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("pgsqlRepository.InitFactory: %w", err)
	}

	rs := &repositories{
		chats:    factory.NewChattRepository(),
		users:    factory.NewUserrRepository(),
		sessions: factory.NewSessionnRepository(),
	}

	closer := func() {
		if err := factory.Close(); err != nil {
			logrus.Error("Закрыть соединение с pgsql: factory.Close: " + err.Error())
		}
	}

	return rs, closer, nil
}
