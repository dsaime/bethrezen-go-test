package createNews

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"newsapi/internal/domain/newsAgr"
)

func upsertReturnsExpected(repo *newsAgr.MockRepository, id int) {
	repo.On("Upsert", mock.Anything).Return(id, nil)
}

func newUsecase(t *testing.T, setupMockRepo func(*newsAgr.MockRepository)) *CreateNewsUsecase {
	repo := newsAgr.NewMockRepository(t)
	if setupMockRepo != nil {
		setupMockRepo(repo)
	}
	return &CreateNewsUsecase{
		Repo: repo,
	}
}

// Test_CreateNews тестирует создание новости
func Test_CreateNews(t *testing.T) {
	t.Run("выходящие совпадают с заданными", func(t *testing.T) {
		const expectedID = 53
		// Настройка мока
		usecase := newUsecase(t, func(repo *newsAgr.MockRepository) {
			upsertReturnsExpected(repo, expectedID)
		})
		// Создать Новость
		in := In{
			Title:      "Some title",
			Content:    "some content",
			Categories: []int{1, 2, 3},
		}
		out, err := usecase.CreateNews(in)
		require.NoError(t, err)
		// Сравнить результат с входящими значениями
		assert.Equal(t, in.Title, out.News.Title)
		assert.Equal(t, in.Content, out.News.Content)
		assert.Equal(t, in.Categories, out.News.Categories)
		// Тот, что вернет репозиторий
		assert.Equal(t, expectedID, out.News.ID)
	})

	t.Run("при невалидных значениях вернет ошибки", func(t *testing.T) {
		// Настройка мока
		usecase := newUsecase(t, nil)
		// Создать новость
		out, err := usecase.CreateNews(In{
			Title:      "",
			Content:    "",
			Categories: []int{1, 1},
		})
		assert.Error(t, err)
		assert.Zero(t, out)
	})
}
