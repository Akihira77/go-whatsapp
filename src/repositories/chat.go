package repositories

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/Akihira77/go_whatsapp/src/store"
	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/Akihira77/go_whatsapp/src/utils"
	"github.com/oklog/ulid/v2"
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
		Preload("Sender").
		Preload("Files").
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
		Preload("Sender").
		Preload("Files").
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
		WithContext(ctx).
		Raw(
			`
            SELECT
                (CASE WHEN messages.group_id IS NULL THEN users.id ELSE NULL END) AS sender_id,
                (CASE WHEN messages.group_id IS NULL THEN users.first_name || ' ' || users.last_name ELSE NULL END) AS user_name,
                (CASE WHEN messages.group_id IS NULL THEN users.image_url ELSE NULL END) AS user_profile,
                (CASE WHEN messages.group_id IS NULL THEN users.status ELSE NULL END) AS user_status,
                COUNT(CASE WHEN messages.is_read = 0 AND messages.sender_id = users.id AND messages.receiver_id IS NOT NULL THEN 1 END) AS unread_peer_chat,
                groups.id AS group_id,
                groups.name AS group_name,
                groups.group_profile,
                COUNT(CASE WHEN messages.is_read = 0 AND messages.group_id = groups.id AND messages.receiver_id IS NULL THEN 1 END) AS unread_group_chat,
                messages.content,
                messages.created_at AS created_at
            FROM messages
            LEFT JOIN users ON users.id = messages.sender_id
            LEFT JOIN groups ON groups.id = messages.group_id
            LEFT JOIN user_groups ON user_groups.group_id = messages.group_id
                WHERE
                    ((users.first_name || ' ' || users.last_name) LIKE ? OR groups.name LIKE ?) AND
                    messages.sender_id <> ? AND 
                    (messages.receiver_id = ? OR user_groups.user_id = ?)
                GROUP BY users.id, groups.id
            UNION
            SELECT
                NULL AS sender_id,
                NULL AS user_name,
                NULL AS user_profile,
                NULL AS user_status,
                0 AS unread_peer_chat,
                groups.id AS group_id,
                groups.name AS group_name,
                groups.group_profile,
                0 AS unread_group_chat,
                NULL,
                groups.created_at AS created_at
            FROM groups
            LEFT JOIN user_groups ON user_groups.group_id = groups.id
            LEFT JOIN users ON users.id = user_groups.user_id
                WHERE
                    groups.name LIKE ? AND
                    groups.id NOT IN (
                        SELECT DISTINCT group_id FROM messages
                        WHERE group_id IS NOT NULL
                    ) AND
                    user_groups.user_id = ?
                GROUP BY sender_id, group_id
                ORDER BY created_at DESC
        `,
			userName,
			groupName,
			myUserId,
			myUserId,
			myUserId,
			groupName,
			myUserId,
		).
		Find(&chatHistories)

	return chatHistories, res.Error
}

func (cr *ChatRepository) AddMessage(ctx context.Context, data types.Message, fileHeaders []*multipart.FileHeader) (*types.Message, error) {
	tx := cr.store.DB.Begin().Debug().WithContext(ctx)

	res := tx.
		Model(&types.Message{}).
		WithContext(ctx).
		Create(&data)

	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	type result struct {
		Data     *bytes.Buffer
		FileName string
		FileType string
		Error    error
	}

	if len(fileHeaders) > 0 {
		results := make(chan result, len(fileHeaders))
		var wg sync.WaitGroup
		var fileList []types.File
		for _, fileHeader := range fileHeaders {
			wg.Add(1)
			go func() {
				defer wg.Done()

				buf, err := utils.ReadFile(fileHeader)
				results <- result{
					Error:    err,
					FileName: fileHeader.Filename,
					FileType: http.DetectContentType(buf.Bytes()),
					Data:     buf,
				}
			}()
		}
		go func() {
			wg.Wait()
			close(results)
		}()

		for result := range results {
			if result.Error != nil {
				tx.Rollback()
				return nil, result.Error
			}

			fileList = append(fileList, types.File{
				ID:        ulid.Make().String(),
				MessageID: data.ID,
				Name:      result.FileName,
				Type:      result.FileType,
				Data:      result.Data.Bytes(),
			})
		}

		res = tx.
			Model(&types.File{}).
			Create(&fileList)
		if res.Error != nil {
			tx.Rollback()
			return nil, res.Error
		}

		data.Files = fileList
	}

	res = tx.Commit()
	return &data, res.Error
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

func (cr *ChatRepository) MarkMessagesAsRead(ctx context.Context, senderId string, receiverId *string, groupId *string) error {
	res := cr.
		store.
		DB.
		Debug().
		WithContext(ctx).
		Model(&types.Message{}).
		Where("(sender_id = ? AND receiver_id IS ? AND group_id IS ?) AND is_read = false", senderId, receiverId, groupId).
		Update("is_read", true)

	return res.Error
}

func (cr *ChatRepository) FindFileInsideChat(ctx context.Context, messageId, fileId string) (*types.File, error) {
	var f types.File
	res := cr.
		store.
		DB.
		Debug().
		WithContext(ctx).
		Model(&types.File{}).
		Where("message_id = ? AND id = ?", messageId, fileId).
		First(&f)

	return &f, res.Error
}

func (cr *ChatRepository) DeleteFile(ctx context.Context, messageId, fileId string) error {
	res := cr.
		store.
		DB.
		Debug().
		WithContext(ctx).
		Model(&types.File{}).
		Where("message_id = ? AND id = ?", messageId, fileId).
		Delete(&types.File{})

	return res.Error
}
