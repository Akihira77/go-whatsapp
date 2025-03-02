package types

import "time"

// INFO: TABLE MODELS
type Message struct {
	ID         string    `json:"id" gorm:"not null;primaryKey"`
	Content    string    `json:"content" gorm:"not null"`
	Files      []File    `json:"files,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsDeleted  bool      `json:"isDeleted" gorm:"not null"`
	IsEdited   bool      `json:"isEdited" gorm:"not null"`
	SenderID   string    `json:"senderId" gorm:"not null"`
	Sender     *User     `json:"sender,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReceiverID *string   `json:"receivedId" gorm:"null"`
	Receiver   *User     `json:"receiver,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	GroupID    *string   `json:"groupId,omitempty" gorm:"null"`
	Group      *Group    `json:"group,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsRead     bool      `json:"isRead" gorm:"not null"`
	CreatedAt  time.Time `json:"createdAt" gorm:"not null"`
}

type File struct {
	ID        string `gorm:"not null;primaryKey"`
	MessageID string `gorm:"not null"`
	Data      []byte `gorm:"type:blob"`
}

// INFO: Data Transfer Object
type CreateMessage struct {
	Content    string  `json:"content,omitempty" form:"content" validate:"max=999"`
	SenderID   string  `json:"senderId" validate:"required"`
	ReceiverID *string `json:"receiverId,omitempty"`
	GroupID    *string `json:"groupId,omitempty"`
}
