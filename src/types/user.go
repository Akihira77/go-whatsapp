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
	ID        string     `json:"id" gorm:"primaryKey"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	ImageUrl  []byte     `json:"imageUrl"`
	Status    UserStatus `json:"userStatus"`
	CreatedAt time.Time  `json:"createdAt"`
}

type UserContact struct {
	UserOneID string    `json:"userOneId" gorm:"primaryKey"`
	UserOne   User      `json:"userOne"`
	UserTwoID string    `json:"userTwoId" gorm:"primaryKey"`
	UserTwo   User      `json:"userTwo"`
	CreatedAt time.Time `json:"createdAt"`
}

// INFO: Data Transfer Object
type Signup struct {
	FirstName string `json:"firstName" form:"firstName" validate:"required"`
	LastName  string `json:"lastName" form:"lastName" validate:"required"`
	Email     string `json:"email" form:"email" validate:"required"`
	Password  string `json:"password" form:"password" validate:"required"`
}

type Signin struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
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
	Search string `json:"search"`
	Pagination
}

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
