package app

import (
	"newsapi/internal/controller/http2"
	mysqlRepository "newsapi/internal/repository/mysql_repository"
)

type Config struct {
	Mysql      mysqlRepository.Config
	Http2      http2.Config
	LogLevel   string
	AuthTokens []string
}
