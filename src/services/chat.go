package services

import (
	"context"
	"fmt"
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

func (cs *ChatService) AddMessage(ctx context.Context, data *types.CreateMessage) (*types.Message, error) {
	newMsg := types.Message{
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

	msg, err := cs.
		chatRepository.
		AddMessage(ctx, newMsg, data.Files)
	if err != nil {
		slog.Error("Add message",
			"error", err,
		)

		return &types.Message{}, fmt.Errorf("Sending chat failed")
	}

	for i := range msg.Files {
		msg.Files[i] = types.File{
			ID:        msg.Files[i].ID,
			MessageID: msg.ID,
			Name:      msg.Files[i].Name,
			Type:      msg.Files[i].Type,
		}
	}

	return msg, nil
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

func (cs *ChatService) MarkMessagesAsRead(ctx context.Context, senderId string, receiverId *string, groupId *string) ([]types.Message, error) {
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

	if groupId != nil {
		return cs.GetMessagesInsideGroup(ctx, *groupId)
	}

	return cs.GetMessages(ctx, senderId, *receiverId)
}

func (cs *ChatService) FindFileInsideChat(ctx context.Context, messageId, fileId string) (*types.File, error) {
	f, err := cs.chatRepository.FindFileInsideChat(ctx, messageId, fileId)
	if err != nil {
		slog.Error("Failed retrieving file",
			"messageId", messageId,
			"fileId", fileId,
			"err", err,
		)
		return nil, err
	}

	return f, nil
}
