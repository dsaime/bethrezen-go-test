package updateNews

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"newsapi/internal/domain/newsAgr"
	mockNewsAgr "newsapi/internal/domain/newsAgr/mocks"
)

func upsertReturnsSameID(repo *mockNewsAgr.Repository) *mock.Call {
	return repo.On("Upsert", mock.Anything).
		Return(func(n newsAgr.News) (int, error) { return n.ID, nil })
}

func inTxReturnsMock(repo *mockNewsAgr.Repository) {
	repo.On("InTransaction", mock.Anything).
		Return(func(fn func(newsAgr.Repository) error) error {
			return fn(repo)
		})
}

func findReturnsExpected(repo *mockNewsAgr.Repository, news newsAgr.News) {
	repo.On("Find", mock.Anything).Return(news, nil)
}

func newUsecase(t *testing.T, setupMockRepo func(*mockNewsAgr.Repository)) *UpdateNewsUsecase {
	repo := mockNewsAgr.NewRepository(t)
	if setupMockRepo != nil {
		setupMockRepo(repo)
	}
	return &UpdateNewsUsecase{
		Repo: repo,
	}
}

// Test_UpdateNews тестирует обновление новости
func Test_UpdateNews(t *testing.T) {
	t.Run("выходящие совпадают с заданными", func(t *testing.T) {
		initialNews := newsAgr.News{
			ID:         42,
			Title:      "Some title",
			Content:    "some content",
			Categories: []int{1, 2, 3},
		}
		// Настройка мока
		usecase := newUsecase(t, func(repo *mockNewsAgr.Repository) {
			findReturnsExpected(repo, initialNews)
			inTxReturnsMock(repo)
			// Обновление 3х полей
			upsertReturnsSameID(repo).Times(3)
		})
		// Обновить новость
		in := In{
			ID:         initialNews.ID,
			Title:      "New title",
			Content:    "new content",
			Categories: []int{4, 5, 6},
		}
		out, err := usecase.UpdateNews(in)
		require.NoError(t, err)
		// Сравнить результат с входящими значениями
		assert.Equal(t, in.ID, out.News.ID)
		assert.Equal(t, in.Title, out.News.Title)
		assert.Equal(t, in.Content, out.News.Content)
		assert.Equal(t, in.Categories, out.News.Categories)
	})

	t.Run("с некорректным id возвращается ошибка", func(t *testing.T) {
		// Настройка мока
		usecase := newUsecase(t, nil)
		// Обновить новость
		_, err := usecase.UpdateNews(In{})
		assert.ErrorIs(t, err, ErrInvalidNewsID)
	})

	t.Run("возвращается ошибка если не указаны новые значения", func(t *testing.T) {
		// Настройка мока
		usecase := newUsecase(t, nil)
		// Нечего обновлять
		_, err := usecase.UpdateNews(In{ID: 31})
		assert.ErrorIs(t, err, ErrNothingToUpdate)
	})

}
