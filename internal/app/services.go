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

func initUsecases(rr *repositories, aa *adapters) usecasesBase {
	return usecasesBase{
		FindSessionsUsecase: &findSession.FindSessionsUsecase{
			Repo: rr.sessions,
		},
		AcceptInvitationUsecase: &acceptInvitation.AcceptInvitationUsecase{
			Repo:          rr.newsRepo,
			EventConsumer: aa.eventBus,
		},
		CancelInvitationUsecase: &cancelInvitation.CancelInvitationUsecase{
			Repo:          rr.newsRepo,
			EventConsumer: aa.eventBus,
		},
		ChatInvitationsUsecase: &chatInvitations.ChatInvitationsUsecase{
			Repo: rr.newsRepo,
		},
		ChatMembersUsecase: &chatMembers.ChatMembersUsecase{
			Repo: rr.newsRepo,
		},
		CreateChatUsecase: &createChat.CreateChatUsecase{
			Repo:          rr.newsRepo,
			EventConsumer: aa.eventBus,
		},
		DeleteMemberUsecase: &deleteMember.DeleteMemberUsecase{
			Repo:          rr.newsRepo,
			EventConsumer: aa.eventBus,
		},
		LeaveChatUsecase: &leaveChat.LeaveChatUsecase{
			Repo:          rr.newsRepo,
			EventConsumer: aa.eventBus,
		},
		MyChatsUsecase: &myChats.MyChatsUsecase{
			Repo: rr.newsRepo,
		},
		ReceivedInvitationsUsecase: &receivedInvitations.ReceivedInvitationsUsecase{
			Repo: rr.newsRepo,
		},
		SendInvitationUsecase: &sendInvitation.SendInvitationUsecase{
			Repo:          rr.newsRepo,
			EventConsumer: aa.eventBus,
		},
		UpdateNameUsecase: &updateName.UpdateNameUsecase{
			Repo:          rr.newsRepo,
			EventConsumer: aa.eventBus,
		},
		BasicAuthRegistrationUsecase: &basicAuthRegistration.BasicAuthRegistrationUsecase{
			Repo:         rr.users,
			SessionsRepo: rr.sessions,
		},
		BasicAuthLoginUsecase: &basicAuthLogin.BasicAuthLoginUsecase{
			Repo:         rr.users,
			SessionsRepo: rr.sessions,
		},
		OauthAuthorizeUsecase: &oauthAuthorize.OauthAuthorizeUsecase{
			Providers: aa.oauthProviders,
		},
		OauthCompleteUsecase: &oauthComplete.OauthCompleteUsecase{
			Repo:         rr.users,
			Providers:    aa.oauthProviders,
			SessionsRepo: rr.sessions,
		},
	}
}
