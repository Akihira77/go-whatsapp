package types

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserStatus string

const (
	ONLINE  UserStatus = "ONLINE"
	OFFLINE UserStatus = "OFFLINE"
)

// INFO: TABLE MODELS
type User struct {
	ID         string     `json:"id" gorm:"not null;primaryKey"`
	FirstName  string     `json:"firstName" gorm:"not null;index:user_name"`
	LastName   string     `json:"lastName" gorm:"not null;index:user_name"`
	Email      string     `json:"email" gorm:"not null;uniqueindex"`
	Password   string     `json:"password" gorm:"not null"`
	ImageUrl   []byte     `json:"imageUrl" gorm:""`
	Status     UserStatus `json:"userStatus" gorm:"not null"`
	LastOnline time.Time  `json:"lastOnline"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"not null"`
}

type UserContact struct {
	UserOneID string    `json:"userOneId" gorm:"not null;primaryKey"`
	UserOne   User      `json:"userOne"`
	UserTwoID string    `json:"userTwoId" gorm:"not null;primaryKey"`
	UserTwo   User      `json:"userTwo"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
}

type Group struct {
	ID           string      `json:"id" gorm:"not null;primaryKey"`
	Name         string      `json:"name" gorm:"not null;index"`
	Description  string      `json:"description"`
	UserCount    int         `json:"userCount" gorm:"not null"`
	Member       []UserGroup `json:"member,omitempty" gorm:"foreignKey:GroupID;references:ID"`
	CreatorID    string      `json:"creatorId" gorm:"not null"`
	Creator      *User       `json:"creator,omitempty"`
	GroupProfile []byte      `json:"groupProfile,omitempty"`
	Messages     []Message   `json:"messages"`
	CreatedAt    time.Time   `json:"createdAt" gorm:"not null"`
}

type UserGroup struct {
	UserID  string `json:"userId" gorm:"not null;primaryKey"`
	User    User   `json:"user,omitempty"`
	GroupID string `json:"groupId" gorm:"not null;primaryKey"`
	Group   Group  `json:"group,omitempty"`
}

// INFO: Data Transfer Object
type ChatDto struct {
	SenderID        string     `json:"senderId"`
	UserName        string     `json:"userName"`
	UserProfile     []byte     `json:"userProfile"`
	UserStatus      UserStatus `json:"userStatus"`
	UnreadPeerChat  int        `json:"unreadPeerChat"`
	GroupID         string     `json:"groupId"`
	GroupName       string     `json:"groupName"`
	GroupProfile    []byte     `json:"groupProfile"`
	UnreadGroupChat int        `json:"unreadGroupChat"`
}

type Signup struct {
	FirstName string `json:"firstName" form:"firstName" validate:"required"`
	LastName  string `json:"lastName" form:"lastName" validate:"required"`
	Email     string `json:"email" form:"email" validate:"required,email"`
	Password  string `json:"password" form:"password" validate:"required,min=6,max=16"`
}

type Signin struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=16"`
}

type UpdateUser struct {
	FirstName string `json:"firstName" form:"firstName"`
	LastName  string `json:"lastName" form:"lastName"`
}

type UpdatePassword struct {
	OldPassword string `json:"oldPassword" form:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" form:"newPassword" validate:"required"`
}

type UserInfo struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	ImageUrl string `json:"imageUrl"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UserQuerySearch struct {
	Search string `json:"search" form:"search"`
	Pagination
}

type Pagination struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type UserMessage struct {
	UserID      string     `json:"userId"`
	FullName    string     `json:"fullName"`
	ImageUrl    []byte     `json:"imageUrl"`
	Status      UserStatus `json:"userStatus"`
	MessageID   string     `json:"messageId"`
	LastMessage string     `json:"lastMessage"`
	Read        bool       `json:"read"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type CreateGroup struct {
	Creator      *User  `json:"creator"`
	Name         string `json:"name" validate:"required"`
	Member       string `json:"member" validate:"required"`
	Description  string `json:"description,omitempty"`
	GroupProfile string `json:"groupProfile,omitempty"`
}

type EditGroup struct {
	EditName        bool   `json:"editName"`
	Name            string `json:"name,omitempty"`
	EditDescription bool   `json:"editDescription"`
	Description     string `json:"description,omitempty"`
}
