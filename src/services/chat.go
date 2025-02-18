package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Akihira77/go_whatsapp/src/repositories"
	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/oklog/ulid/v2"
)

type ChatService struct {
	chatRepository *repositories.ChatRepository
}

func NewChatService(chatRepository *repositories.ChatRepository) *ChatService {
	return &ChatService{
		chatRepository: chatRepository,
	}
}

func (cs *ChatService) GetMessages(ctx context.Context, userOneId, userTwoId string) ([]types.Message, error) {
	msgs, err := cs.
		chatRepository.
		GetMessages(ctx, [2]string{userOneId, userTwoId})
	if err != nil {
		slog.Error("Retrieving messages inside chat room",
			"error", err,
		)
	}

	return msgs, err
}

func (cs *ChatService) GetMessagesInsideGroup(ctx context.Context, groupId string) ([]types.Message, error) {
	msgs, err := cs.
		chatRepository.
		GetMessagesInsideGroup(ctx, groupId)
	if err != nil {
		slog.Error("Retrieving messages inside group",
			"error", err,
		)
	}

	return msgs, err
}

func (cs *ChatService) SearchChat(ctx context.Context, myUserId, userName, groupName string) ([]types.ChatDto, error) {
	chatHistories, err := cs.
		chatRepository.
		SearchChat(ctx, myUserId, groupName, userName)
	if err != nil {
		slog.Error("Retrieving last message on sidebar",
			"error", err,
		)
	}

	return chatHistories, err
}

func (cs *ChatService) AddMessage(ctx context.Context, data *types.CreateMessage) (types.Message, error) {
	msg := types.Message{
		ID:         ulid.Make().String(),
		Content:    data.Content,
		SenderID:   data.SenderID,
		ReceiverID: data.ReceiverID,
		GroupID:    data.GroupID,
		IsEdited:   false,
		IsDeleted:  false,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}

	err := cs.
		chatRepository.
		AddMessage(ctx, msg)
	if err != nil {
		slog.Error("Add message",
			"error", err,
		)

		return types.Message{}, err
	}

	return msg, err
}

func (cs *ChatService) EditMessage(ctx context.Context, data *types.Message) (types.Message, error) {
	data.IsEdited = true

	err := cs.
		chatRepository.
		EditMessage(ctx, *data)
	if err != nil {
		slog.Error("Edit message",
			"error", err,
		)

		return types.Message{}, err
	}

	return *data, err
}

func (cs *ChatService) SoftDeleteMessage(ctx context.Context, data *types.Message) (types.Message, error) {
	data.IsDeleted = true

	err := cs.
		chatRepository.
		EditMessage(ctx, *data)
	if err != nil {
		slog.Error("Soft delete message",
			"error", err,
		)

		return types.Message{}, err
	}

	return *data, err
}

func (cs *ChatService) MarkMessagesAsRead(ctx context.Context, senderId, receiverId, groupId string) ([]types.Message, error) {
	err := cs.
		chatRepository.
		MarkMessagesAsRead(ctx, senderId, receiverId, groupId)
	if err != nil {
		slog.Error("Marking messages as read failed",
			"sender", senderId,
			"receiver", receiverId,
		)
		return []types.Message{}, err
	}

	if groupId != "" {
		return cs.GetMessagesInsideGroup(ctx, groupId)
	}

	return cs.GetMessages(ctx, senderId, receiverId)
}
