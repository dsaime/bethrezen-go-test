package newsList

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"newsapi/internal/domain/newsAgr"
)

func listReturnsExpected(repo *newsAgr.MockRepository, news []newsAgr.News) *mock.Call {
	return repo.On("List", mock.Anything).
		Return(news, nil)
}

func newUsecase(t *testing.T, setupMockRepo func(*newsAgr.MockRepository)) *NewsListUsecase {
	repo := newsAgr.NewMockRepository(t)
	if setupMockRepo != nil {
		setupMockRepo(repo)
	}
	return &NewsListUsecase{
		Repo: repo,
	}
}

// Test_UpdateNews тестирует обновление новости
func Test_UpdateNews(t *testing.T) {
	t.Run("возвращает все из репозитория", func(t *testing.T) {
		expected := []newsAgr.News{
			{ID: 31, Title: "Title", Content: "Content"},
			{ID: 32, Title: "Title1", Content: "Content2"},
			{ID: 34, Title: "Title2", Content: "Content5"},
		}
		// Настройка мока
		usecase := newUsecase(t, func(repo *newsAgr.MockRepository) {
			listReturnsExpected(repo, expected)
		})
		// Обновить новость
		out, err := usecase.NewsList(In{})
		require.NoError(t, err)
		// Сравнить результат с входящими значениями
		assert.Equal(t, expected, out.News)
	})
}
