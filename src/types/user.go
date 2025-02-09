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
	ID         string     `json:"id" gorm:"primaryKey"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	ImageUrl   []byte     `json:"imageUrl"`
	Status     UserStatus `json:"userStatus"`
	LastOnline time.Time  `json:"lastOnline"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type UserContact struct {
	UserOneID string    `json:"userOneId" gorm:"primaryKey"`
	UserOne   User      `json:"userOne"`
	UserTwoID string    `json:"userTwoId" gorm:"primaryKey"`
	UserTwo   User      `json:"userTwo"`
	CreatedAt time.Time `json:"createdAt"`
}

type Group struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"index"`
	UserCount int       `json:"userCount"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserGroup struct {
	UserID  string `json:"userId" gorm:"primaryKey"`
	User    User   `json:"user,omitempty"`
	GroupID string `json:"groupId" gorm:"primaryKey"`
	Group   Group  `json:"group,omitempty"`
}

// INFO: Data Transfer Object
type UserDto struct {
	ID              string     `json:"id" gorm:"primaryKey"`
	FullName        string     `json:"fullName"`
	ImageUrl        []byte     `json:"imageUrl"`
	Status          UserStatus `json:"userStatus"`
	UnreadChatCount int        `json:"unreadChatCount"`
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
