package app

import (
	"newsapi/internal/controller/http2"
	mysqlRepository "newsapi/internal/repository/mysql_repository"
)

// Config представляет собой композицию из конфигураций зависимостей приложения
type Config struct {
	Mysql      mysqlRepository.Config // Настройка mysql
	Http2      http2.Config           // Настройка сервера http
	LogLevel   string                 // Уровень логирования
	AuthTokens []string               // Список токенов авторизации
}
