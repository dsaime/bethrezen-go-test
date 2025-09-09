package newsList

import (
	"testing"

	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"

	serviceSuite "newsapi/internal/usecases/suite"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

// Test_Chats_MyChats тестирует запрос список чатов в которых участвует пользователь
func (suite *testSuite) Test_Chats_MyChats() {
	usecase := &NewsListUsecase{
		Repo: suite.RR.Chats,
	}

	suite.Run("пользователь может запрашивать только свой чат", func() {
		input := In{
			SubjectID: uuid.New(),
			UserID:    uuid.New(),
		}
		out, err := usecase.NewsList(input)
		suite.ErrorIs(err, ErrUnauthorizedChatsView)
		suite.Empty(out)
	})

	suite.Run("пустой список из пустого репозитория", func() {
		input := suite.newUserChatsInput(uuid.New())
		out, err := usecase.NewsList(input)
		suite.NoError(err)
		suite.Empty(out)
	})

	suite.Run("пустой список если у пользователя нет чатов", func() {
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := suite.UpsertChat(suite.RndChat())
			suite.AddRndParticipant(&chat)
		}
		input := suite.newUserChatsInput(uuid.New())
		out, err := usecase.NewsList(input)
		suite.NoError(err)
		suite.Empty(out)
	})

	suite.Run("у пользователя может быть несколько чатов", func() {
		userID := uuid.New()
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := suite.RndChat()
			p := suite.NewParticipant(userID)
			suite.AddParticipant(&chat, p)
			suite.UpsertChat(chat)
		}
		input := suite.newUserChatsInput(userID)
		out, err := usecase.NewsList(input)
		suite.NoError(err)
		suite.Len(out.News, chatsAllCount)
	})
}

func (suite *testSuite) newUserChatsInput(userID uuid.UUID) In {
	return In{
		SubjectID: userID,
		UserID:    userID,
	}
}
