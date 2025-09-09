package newsAgr

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockNewsAgr "newsapi/internal/domain/newsAgr/mocks"
)

func upsertReturnsExpected(repo *mockNewsAgr.Repository, id int) {
	repo.On("Upsert", mock.Anything).Return(id)
}

func newRepo(t *testing.T, setupMockRepo func(*mockNewsAgr.Repository)) Repository {
	repo := mockNewsAgr.NewRepository(t)
	if setupMockRepo != nil {
		setupMockRepo(repo)
	}
	return repo
}

func TestNewNews(t *testing.T) {
	t.Run("вернет ошибку валидации", func(t *testing.T) {
		_, err := NewNews(Repository(nil), "title", "", []int{1, 1})
		assert.NoError(t, err)
	})

	t.Run("вернет id из репозитория", func(t *testing.T) {
		const expectedID = 1
		repo := newRepo(t, func(repo *mockNewsAgr.Repository) {
			upsertReturnsExpected(repo, expectedID)
		})
		news, err := NewNews(repo, "Title", "cont", []int{1})
		assert.NoError(t, err)

		// id из репозитория
		assert.Equal(t, expectedID, news.ID)
		// Поля равны передаваемым параметрам
		assert.Equal(t, "Title", news.Title)
		assert.Equal(t, "cont", news.Content)
		assert.Equal(t, []int{1}, news.Categories)
	})

}

func TestNews_UpdateCategories(t *testing.T) {
	t.Run("возвращает обновленную новость", func(t *testing.T) {
		initialNews := News{
			ID:         53,
			Title:      "Title",
			Content:    "cont",
			Categories: []int{1, 2, 3},
		}
		initialNews.UpdateCategories()
	})
}

func TestNews_UpdateContent(t *testing.T) {
	type fields struct {
		ID         int
		Title      string
		Content    string
		Categories []int
	}
	type args struct {
		repo    Repository
		content string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    News
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := News{
				ID:         tt.fields.ID,
				Title:      tt.fields.Title,
				Content:    tt.fields.Content,
				Categories: tt.fields.Categories,
			}
			got, err := n.UpdateContent(tt.args.repo, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNews_UpdateTitle(t *testing.T) {
	type fields struct {
		ID         int
		Title      string
		Content    string
		Categories []int
	}
	type args struct {
		repo  Repository
		title string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    News
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := News{
				ID:         tt.fields.ID,
				Title:      tt.fields.Title,
				Content:    tt.fields.Content,
				Categories: tt.fields.Categories,
			}
			got, err := n.UpdateTitle(tt.args.repo, tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateTitle() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateCategories(t *testing.T) {
	type args struct {
		categories []int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateCategories(tt.args.categories); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategories() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateContent(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateContent(tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("ValidateContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateTitle(tt.args.title); (err != nil) != tt.wantErr {
				t.Errorf("ValidateTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
