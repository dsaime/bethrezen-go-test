package app

import (
	oauthProvider "newsapi/internal/adapter/oauth_provider"
	"newsapi/internal/controller/http2"
	pgsqlRepository "newsapi/internal/repository/pgsql_repository"
)

type Config struct {
	Pgsql       pgsqlRepository.Config
	Http2       http2.Config
	LogLevel    string
	OauthGoogle oauthProvider.GoogleConfig
	OauthGithub oauthProvider.GithubConfig
}
