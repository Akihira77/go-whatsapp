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

func (cr *ChatRepository) GetMessagesInsideGroup(ctx context.Context, groupId string) ([]types.Message, error) {
	var msgs []types.Message

	res := cr.
		store.
		DB.
		Debug().
		Model(&types.Message{}).
		WithContext(ctx).
		Where("group_id", groupId).
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

func (cr *ChatRepository) SearchChat(ctx context.Context, myUserId, groupName, userName string) ([]types.ChatDto, error) {
	var chatHistories []types.ChatDto

	res := cr.
		store.
		DB.
		Debug().
		Model(&types.Group{}).
		WithContext(ctx).
		Joins("LEFT JOIN user_groups ON user_groups.group_id = groups.id").
		Joins("LEFT JOIN messages ON messages.group_id = groups.id").
		Joins("LEFT JOIN users ON users.id = messages.sender_id").
		Select(
			"(CASE WHEN groups.name IS NULL THEN users.id ELSE NULL END) AS sender_id",
			"(CASE WHEN groups.name IS NULL THEN users.first_name || ' ' || users.last_name ELSE NULL END) AS user_name",
			"(CASE WHEN groups.name IS NULL THEN users.image_url ELSE NULL END) AS user_profile",
			"COUNT(CASE WHEN messages.is_read = 0 AND messages.sender_id = users.id AND messages.receiver_id IS NOT NULL THEN 1 END) AS unread_peer_chat",
			"groups.id AS group_id",
			"groups.name AS group_name",
			"groups.group_profile",
			"COUNT(CASE WHEN messages.is_read = 0 AND messages.group_id = groups.id AND messages.group_id IS NOT NULL THEN 1 END) AS unread_group_chat",
		).
		Where(
			"(messages.receiver_id = ? AND (users.first_name || ' ' || users.last_name) LIKE ?) OR groups.name LIKE ?",
			myUserId,
			userName,
			groupName,
		).
		Group("users.id, groups.id").
		Order("MAX(messages.created_at) DESC").
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

func (cr *ChatRepository) MarkMessagesAsRead(ctx context.Context, senderId, receiverId, groupId string) error {
	res := cr.
		store.
		DB.
		Model(&types.Message{}).
		WithContext(ctx).
		Where("(sender_id = ? AND receiver_id = ? AND group_id = ?) AND is_read = false", senderId, receiverId, groupId).
		Update("is_read", true)

	return res.Error
}
