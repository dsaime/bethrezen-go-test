package mysqlRepository

import (
	"errors"

	"newsapi/internal/domain/newsAgr"
)

func (suite *Suite) TestNewsRepository() {
	suite.Run("Find/Найдет запись по id", func() {
		// Вставить записи
		inNews := []newsAgr.News{
			{ID: 1, Title: "A", Content: "A", Categories: []int{1, 2, 3}},
		}
		suite.insertNews(inNews)

		// Поиск
		news, err := suite.RR.News.Find(newsAgr.Filter{ID: 1})
		suite.NoError(err)
		suite.Equal(inNews[0], news)
	})

	suite.Run("Find/Вернет ошибку если не найдет", func() {
		// Вставить записи
		inNews := []newsAgr.News{
			{ID: 1, Title: "A", Content: "A", Categories: []int{1, 2, 3}},
		}
		suite.insertNews(inNews)

		// Поиск
		_, err := suite.RR.News.Find(newsAgr.Filter{ID: 2})
		suite.ErrorIs(err, newsAgr.ErrNewsNotFound)
	})

	suite.Run("InTransaction/Не сохранится в БД после отмены", func() {
		var id int
		err := suite.RR.News.InTransaction(func(txRepo newsAgr.Repository) (err error) {
			// Сделать вставку
			id, err = txRepo.Upsert(newsAgr.News{
				Title:      "A",
				Content:    "A",
				Categories: []int{1, 2, 3},
			})
			suite.Require().NoError(err)

			// Проверить запись
			_, err = txRepo.Find(newsAgr.Filter{ID: id})
			suite.Require().NoError(err)

			// Завершить с ошибкой
			return errors.New("foo")
		})
		// Проверить результат выполнения
		suite.Equal("foo", err.Error())
		suite.Equal(1, id)
		// Убедиться, что записи нет
		_, err = suite.RR.News.Find(newsAgr.Filter{ID: id})
		suite.ErrorIs(err, newsAgr.ErrNewsNotFound)
	})

	suite.Run("List/Вернет все записи", func() {
		// Вставить записи
		inNews := []newsAgr.News{
			{Title: "A", Content: "A", Categories: []int{1, 2, 3}},
			{Title: "B", Content: "B", Categories: []int{1, 3}},
			{Title: "C", Content: "C", Categories: []int{3}},
		}
		suite.insertNews(inNews)

		// После появятся ID вставки
		inNews[0].ID = 1
		inNews[1].ID = 2
		inNews[2].ID = 3

		news, err := suite.RR.News.List(newsAgr.Filter{})
		suite.Require().NoError(err)
		suite.Require().Equal(inNews, news)
	})
	suite.Run("Upsert/После вставки можно прочитать из репозитория", func() {
		// Сделать вставку
		id, err := suite.RR.News.Upsert(newsAgr.News{
			Title:      "A",
			Content:    "B",
			Categories: []int{4, 5},
		})
		suite.NoError(err)
		suite.Equal(1, id)

		// Поверить новости
		var savedNews dbNews
		err = suite.factory.db.Get(&savedNews, `SELECT * FROM News WHERE id = ?`, id)
		suite.NoError(err)
		suite.Equal(1, savedNews.ID)
		suite.Equal("A", savedNews.Title)
		suite.Equal("B", savedNews.Content)

		// Проверить категории
		var savedCats []dbCategory
		err = suite.factory.db.Select(&savedCats, `SELECT * FROM NewsCategories`)
		suite.NoError(err)
		suite.Require().Len(savedCats, 2)
		suite.Equal(dbCategory{ID: 4, NewsID: 1}, savedCats[0])
		suite.Equal(dbCategory{ID: 5, NewsID: 1}, savedCats[1])
	})
}

func (suite *Suite) insertNews(inNews []newsAgr.News) {
	for _, news := range inNews {
		res, err := suite.factory.db.Exec(`
				INSERT INTO News (Title, Content)
				VALUES (?, ?)
			`, news.Title, news.Content)
		suite.Require().NoError(err)

		newsID, _ := res.LastInsertId()

		for _, category := range news.Categories {
			_, err = suite.factory.db.Exec(`
				INSERT INTO NewsCategories (NewsId, CategoryId)
				VALUES (?, ?)
			`, newsID, category)
			suite.Require().NoError(err)
		}
	}
}
