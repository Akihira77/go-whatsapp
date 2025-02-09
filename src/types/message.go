package types

import "time"

// INFO: TABLE MODELS
type Message struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	Content    string    `json:"content"`
	Files      []File    `json:"files,omitempty" gorm:"foreignKey:MessageID"`
	IsDeleted  bool      `json:"isDeleted"`
	IsEdited   bool      `json:"isEdited"`
	SenderID   string    `json:"senderId"`
	Sender     *User     `json:"sender,omitempty"`
	ReceiverID string    `json:"receivedId"`
	Receiver   *User     `json:"receiver,omitempty"`
	IsRead     bool      `json:"isRead"`
	CreatedAt  time.Time `json:"createdAt"`
}

type File struct {
	ID        string `gorm:"primaryKey"`
	MessageID string
	Data      []byte `gorm:"type:blob"`
}

// INFO: Data Transfer Object
type CreateMessage struct {
	Content    string `json:"content,omitempty" form:"content" validate:"max=999"`
	SenderID   string `json:"senderId" validate:"required"`
	ReceiverID string `json:"receiverId" validate:"required"`
}
