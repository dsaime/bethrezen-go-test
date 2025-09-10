package app

import (
	createNews "newsapi/internal/usecases/news/create_news"
	newsList "newsapi/internal/usecases/news/news_list"
	updateNews "newsapi/internal/usecases/news/update_news"
)

// usecasesBase содержит инициализированные usecases
type usecasesBase struct {
	// News
	*createNews.CreateNewsUsecase
	*newsList.NewsListUsecase
	*updateNews.UpdateNewsUsecase
}

// initUsecases инициализирует usecases
func initUsecases(rr *repositories) usecasesBase {
	return usecasesBase{
		CreateNewsUsecase: &createNews.CreateNewsUsecase{Repo: rr.newsRepo},
		NewsListUsecase:   &newsList.NewsListUsecase{Repo: rr.newsRepo},
		UpdateNewsUsecase: &updateNews.UpdateNewsUsecase{Repo: rr.newsRepo},
	}
}
