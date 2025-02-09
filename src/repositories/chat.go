package repositories

import (
	"context"

	"github.com/Akihira77/go_whatsapp/src/store"
	"github.com/Akihira77/go_whatsapp/src/types"
)

type ChatRepository struct {
	store *store.Store
}

func NewChatRepository(store *store.Store) *ChatRepository {
	return &ChatRepository{
		store: store,
	}
}

func (cr *ChatRepository) GetMessages(ctx context.Context, userIds [2]string) ([]types.Message, error) {
	var msgs []types.Message

	res := cr.
		store.
		DB.
		Debug().
		Model(&types.Message{}).
		WithContext(ctx).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userIds[0], userIds[1], userIds[1], userIds[0]).
		Find(&msgs)

	return msgs, res.Error
}

func (cr *ChatRepository) SearchMessageInsideChat(ctx context.Context, userIds []string, content string) ([]types.Message, error) {
	var msgs []types.Message

	res := cr.
		store.
		DB.
		Model(&types.Message{}).
		WithContext(ctx).
		Where(
			"((sender_id = ? AND receiver_id = ?) OR (sender_id = ? OR receiver_id = ?)) AND content LIKE ?",
			userIds[0],
			userIds[1],
			userIds[1],
			userIds[0],
			"%"+content+"%",
		).
		Find(&msgs)

	return msgs, res.Error
}

func (cr *ChatRepository) SearchChat(ctx context.Context, myUserId, name string) ([]types.UserDto, error) {
	var chatHistories []types.UserDto

	res := cr.
		store.
		DB.
		Debug().
		Model(&types.Message{}).
		WithContext(ctx).
		Joins("JOIN users ON users.id = messages.sender_id").
		Where(
			"messages.receiver_id = ? AND (users.first_name || ' ' || users.last_name) LIKE ?",
			myUserId,
			"%"+name+"%",
		).
		Order("messages.created_at DESC").
		Select(
			"users.id AS id",
			"users.first_name || ' ' || users.last_name AS full_name",
			"users.image_url",
			"users.status",
			"SUM(CASE WHEN messages.is_read = 0 THEN 1 ELSE 0 END) AS unread_chat_count",
		).
		Group("users.id").
		Find(&chatHistories)

	return chatHistories, res.Error
}

func (cr *ChatRepository) AddMessage(ctx context.Context, data types.Message) error {
	res := cr.
		store.
		DB.
		Model(&types.Message{}).
		WithContext(ctx).
		Create(&data)

	return res.Error
}

func (cr *ChatRepository) EditMessage(ctx context.Context, data types.Message) error {
	res := cr.
		store.
		DB.
		WithContext(ctx).
		Save(&data)

	return res.Error
}

func (cr *ChatRepository) HardDeleteMessage(ctx context.Context, msgId string) error {
	res := cr.
		store.
		DB.
		Model(&types.Message{}).
		WithContext(ctx).
		Where("id = ?", msgId).
		Delete(&types.Message{})

	return res.Error
}

func (cr *ChatRepository) MarkMessagesAsRead(ctx context.Context, senderId, receiverId string) error {
	res := cr.
		store.
		DB.
		Model(&types.Message{}).
		WithContext(ctx).
		Where("(sender_id = ? AND receiver_id = ?) AND is_read = false", senderId, receiverId).
		Update("is_read", true)

	return res.Error
}
