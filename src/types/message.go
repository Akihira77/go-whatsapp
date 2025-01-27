package types

// INFO: TABLE MODELS
type Message struct {
	ID         string `json:"id" gorm:"primaryKey"`
	Content    string `json:"content"`
	Files      []File `json:"files" gorm:"foreignKey:MessageID"`
	IsDeleted  bool   `json:"isDeleted"`
	IsEdited   bool   `json:"isEdited"`
	SenderID   string `json:"senderId"`
	Sender     User   `json:"sender"`
	ReceiverID string `json:"receivedId"`
	Receiver   User   `json:"received"`
	CreatedAt  string `json:"createdAt"`
}

type File struct {
	ID        string `gorm:"primaryKey"`
	MessageID string
	Data      []byte `gorm:"type:blob"`
}
