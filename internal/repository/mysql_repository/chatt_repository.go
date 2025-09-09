package mysqlRepository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nullism/bqb"

	"newsapi/internal/domain/newsAgr"
	sqlxRepo "newsapi/internal/repository/mysql_repository/sqlx_repo"
)

type NewsRepository struct {
	sqlxRepo.SqlxRepo
}

func (r *NewsRepository) Find(filter newsAgr.Filter) (newsAgr.News, error) {
	//TODO implement me
	panic("implement me")
}

func (r *NewsRepository) List(filter newsAgr.Filter) ([]newsAgr.News, error) {
	//TODO implement me
	panic("implement me")
}

func (r *NewsRepository) Upsert(news newsAgr.News) (id int, _ error) {
	//TODO implement me
	panic("implement me")
}

func (r *NewsRepository) List(filter newsAgr.Filter) ([]newsAgr.Chat, error) {
	sel := bqb.New("SELECT c.* FROM chats c")
	where := bqb.Optional("WHERE")

	needJoinParticipants := filter.ParticipantID != uuid.Nil
	if needJoinParticipants {
		sel = sel.Space("LEFT JOIN participants p ON c.id = p.chat_id")
	}
	if filter.ParticipantID != uuid.Nil {
		where = where.And("p.user_id = ?", filter.ParticipantID)
	}

	needJoinInvitations := filter.InvitationID != uuid.Nil || filter.InvitationRecipientID != uuid.Nil
	if needJoinInvitations {
		sel = sel.Space("LEFT JOIN invitations i ON c.id = i.chat_id")
	}
	if filter.InvitationID != uuid.Nil {
		where = where.And("i.id = ?", filter.InvitationID)
	}
	if filter.InvitationRecipientID != uuid.Nil {
		where = where.And("i.recipient_id = ?", filter.InvitationRecipientID)
	}

	if filter.ID != uuid.Nil {
		where = where.And("c.id = ?", filter.ID)
	}

	query, args, err := bqb.New("? ? GROUP BY c.id", sel, where).ToPgsql()
	if err != nil {
		return nil, fmt.Errorf("bqb.ToPgsql: %w", err)
	}

	// Запросить чаты
	var chats []dbChat
	if err := r.DB().Select(&chats, query, args...); err != nil {
		return nil, fmt.Errorf("bqb.Select: %w", err)
	}

	// Если чатов нет, сразу вернуть пустой список
	if len(chats) == 0 {
		return nil, nil
	}

	// Собрать ID найденных чатов
	chatIDs := make([]string, len(chats))
	for i, c := range chats {
		chatIDs[i] = c.ID
	}

	// Найти участников чатов
	var participants []dbParticipant
	if err := r.DB().Select(&participants, `
		SELECT *
		FROM participants
		WHERE chat_id = ANY($1)
	`, pq.Array(chatIDs)); err != nil {
		return nil, fmt.Errorf("r.DB().Select: %w", err)
	}

	// Создать карту, где ключ это ID чата, а значение это список его участников
	participantsMap := make(map[string][]dbParticipant, len(chats))
	for _, p := range participants {
		participantsMap[p.ChatID] = append(participantsMap[p.ChatID], p)
	}

	// Найти приглашения в чат
	var invitations []dbInvitation
	if err := r.DB().Select(&invitations, `
		SELECT *
		FROM invitations
		WHERE chat_id = ANY($1)
	`, pq.Array(chatIDs)); err != nil {
		return nil, fmt.Errorf("r.DB().Select: %w", err)
	}

	// Создать карту, где ключ это ID чата, а значение это список приглашений в него
	invitationsMap := make(map[string][]dbInvitation, len(chats))
	for _, i := range invitations {
		invitationsMap[i.ChatID] = append(invitationsMap[i.ChatID], i)
	}

	return toDomainChats(chats, participantsMap, invitationsMap), nil
}

func (r *NewsRepository) Upsert(chat newsAgr.Chat) error {
	if chat.ID == uuid.Nil {
		return fmt.Errorf("chat ID is required")
	}

	if r.IsTx() {
		return r.upsert(chat)
	} else {
		return r.InTransaction(func(txRepo newsAgr.Repository) error {
			return txRepo.Upsert(chat)
		})
	}
}

func (r *NewsRepository) upsert(chat newsAgr.Chat) error {
	if _, err := r.DB().NamedExec(`
		INSERT INTO chats(id, name, chief_id) 
		VALUES (:id, :name, :chief_id)
		ON CONFLICT (id) DO UPDATE SET
			name=excluded.name,
			chief_id=excluded.chief_id
	`, toDBChat(chat)); err != nil {
		return fmt.Errorf("r.DB().NamedExec: %w", err)
	}

	// Удалить прошлых участников
	if _, err := r.DB().Exec(`
		DELETE FROM participants WHERE chat_id = $1
	`, chat.ID); err != nil {
		return fmt.Errorf("r.DB().Exec: %w", err)
	}

	if len(chat.Participants) > 0 {
		if _, err := r.DB().NamedExec(`
			INSERT INTO participants(chat_id, user_id)
			VALUES (:chat_id, :user_id)
		`, toDBParticipants(chat)); err != nil {
			return fmt.Errorf("r.DB().NamedExec: %w", err)
		}
	}

	// Удалить прошлые приглашения
	if _, err := r.DB().Exec(`
		DELETE FROM invitations WHERE chat_id = $1
	`, chat.ID); err != nil {
		return fmt.Errorf("r.DB().Exec: %w", err)
	}

	if len(chat.Invitations) > 0 {
		if _, err := r.DB().NamedExec(`
		INSERT INTO invitations(id, chat_id, subject_id, recipient_id)
		VALUES (:id, :chat_id, :subject_id, :recipient_id)
	`, toDBInvitations(chat)); err != nil {
			return fmt.Errorf("r.DB().NamedExec: %w", err)
		}
	}

	return nil
}

func (r *NewsRepository) InTransaction(fn func(txRepo newsAgr.Repository) error) error {
	return r.SqlxRepo.InTransaction(func(txSqlxRepo sqlxRepo.SqlxRepo) error {
		return fn(&NewsRepository{SqlxRepo: txSqlxRepo})
	})
}

type dbChat struct {
	ID      string `db:"id"`
	Name    string `db:"name"`
	ChiefID string `db:"chief_id"`
}

func toDBChat(chat newsAgr.Chat) dbChat {
	return dbChat{
		ID:      chat.ID.String(),
		Name:    chat.Name,
		ChiefID: chat.ChiefID.String(),
	}
}

func toDomainChat(
	chat dbChat,
	participants []dbParticipant,
	invitations []dbInvitation,
) newsAgr.Chat {
	return newsAgr.Chat{
		ID:           uuid.MustParse(chat.ID),
		Name:         chat.Name,
		ChiefID:      uuid.MustParse(chat.ChiefID),
		Participants: toDomainParticipants(participants),
		Invitations:  toDomainInvitations(invitations),
	}
}

func toDomainChats(
	chats []dbChat,
	participants map[string][]dbParticipant,
	invitations map[string][]dbInvitation,
) []newsAgr.Chat {
	domainChats := make([]newsAgr.Chat, len(chats))
	for i, chat := range chats {
		domainChats[i] = toDomainChat(chat, participants[chat.ID], invitations[chat.ID])
	}

	return domainChats
}

type dbParticipant struct {
	ChatID string `db:"chat_id"`
	UserID string `db:"user_id"`
}

func toDBParticipants(chat newsAgr.Chat) []dbParticipant {
	dbParticipants := make([]dbParticipant, len(chat.Participants))
	for i, p := range chat.Participants {
		dbParticipants[i] = dbParticipant{
			ChatID: chat.ID.String(),
			UserID: p.UserID.String(),
		}
	}

	return dbParticipants
}
func toDomainParticipants(participants []dbParticipant) []newsAgr.Participant {
	pp := make([]newsAgr.Participant, len(participants))
	for i, p := range participants {
		pp[i] = newsAgr.Participant{
			UserID: uuid.MustParse(p.UserID),
		}
	}

	return pp
}

type dbInvitation struct {
	ID          string `db:"id"`
	ChatID      string `db:"chat_id"`
	SubjectID   string `db:"subject_id"`
	RecipientID string `db:"recipient_id"`
}

func toDBInvitations(chat newsAgr.Chat) []dbInvitation {
	dbInvitations := make([]dbInvitation, len(chat.Invitations))
	for i, inv := range chat.Invitations {
		dbInvitations[i] = dbInvitation{
			ID:          inv.ID.String(),
			ChatID:      chat.ID.String(),
			SubjectID:   inv.SubjectID.String(),
			RecipientID: inv.RecipientID.String(),
		}
	}

	return dbInvitations
}

func toDomainInvitations(invitations []dbInvitation) []newsAgr.Invitation {
	ii := make([]newsAgr.Invitation, len(invitations))
	for i, inv := range invitations {
		ii[i] = newsAgr.Invitation{
			ID:          uuid.MustParse(inv.ID),
			RecipientID: uuid.MustParse(inv.RecipientID),
			SubjectID:   uuid.MustParse(inv.SubjectID),
		}
	}

	return ii
}
