package app

import (
	"github.com/sirupsen/logrus"

	eventsBus "newsapi/internal/adapter/events_bus"
	oauthProvider "newsapi/internal/adapter/oauth_provider"
	"newsapi/internal/usecases/users/oauth"
)

type adapters struct {
	oauthProviders oauth.Providers
	eventBus       *eventsBus.EventsBus
}

func (a *adapters) OauthProviders() oauth.Providers {
	return a.oauthProviders
}

func initAdapters(cfg Config) *adapters {
	oauthProviders := oauth.Providers{}
	if cfg.OauthGoogle != (oauthProvider.GoogleConfig{}) {
		oauthProviders.Add(oauthProvider.NewGoogle(cfg.OauthGoogle))
		logrus.Info("Подключен Oauth провайдер Google")
	}
	if cfg.OauthGithub != (oauthProvider.GithubConfig{}) {
		oauthProviders.Add(oauthProvider.NewGithub(cfg.OauthGithub))
		logrus.Info("Подключен Oauth провайдер Github")
	}

	return &adapters{
		oauthProviders: oauthProviders,
		eventBus:       new(eventsBus.EventsBus),
	}
}
