package app

import (
	createNews "newsapi/internal/usecases/news/create_news"
	newsList "newsapi/internal/usecases/news/news_list"
	updateNews "newsapi/internal/usecases/news/update_news"
)

type usecasesBase struct {
	// News
	*createNews.CreateNewsUsecase
	*newsList.NewsListUsecase
	*updateNews.UpdateNewsUsecase
}

func initUsecases(rr *repositories) usecasesBase {
	return usecasesBase{
		CreateNewsUsecase: &createNews.CreateNewsUsecase{Repo: rr.newsRepo},
		NewsListUsecase:   &newsList.NewsListUsecase{Repo: rr.newsRepo},
		UpdateNewsUsecase: &updateNews.UpdateNewsUsecase{Repo: rr.newsRepo},
	}
}
