package utils

import "github.com/Akihira77/go_whatsapp/src/types"

func GetFullName(user *types.User) string {
	return user.FirstName + " " + user.LastName
}
