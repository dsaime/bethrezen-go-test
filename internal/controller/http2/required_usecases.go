package http2

import "newsapi/internal/controller/http2/register_handler"

// RequiredUsecases определяет интерфейс для доступа к сценариям использования бизнес-логики
type RequiredUsecases interface {
	registerHandler.UsecasesForNewsList
	registerHandler.UsecasesForUpdateNews
	registerHandler.UsecasesForCreateNews
}
