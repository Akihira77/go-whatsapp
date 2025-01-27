package types

import (
	"mime/multipart"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// INFO: TABLE MODELS
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	ImageUrl  []byte    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
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
	FirstName string                `json:"firstName" form:"firstName"`
	LastName  string                `json:"lastName" form:"lastName"`
	Email     string                `json:"email" form:"email"`
	Password  string                `json:"password" form:"password"`
	Image     *multipart.FileHeader `json:"image" form:"image"`
}

type Signin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUser struct {
	FirstName string                `json:"firstName" form:"firstName"`
	LastName  string                `json:"lastName" form:"lastName"`
	Image     *multipart.FileHeader `json:"image" form:"image"`
}

type UpdatePassword struct {
	OldPassword string `json:"oldPassword" form:"oldPassword"`
	NewPassword string `json:"newPassword" form:"newPassword"`
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
