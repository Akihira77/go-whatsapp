package types

import "time"

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
type Signin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	ImageUrl string `json:"imageUrl"`
}
