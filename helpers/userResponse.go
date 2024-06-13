package helper

import (
	"spotify/internal/database"

	"github.com/google/uuid"
)

type CustomUser struct {
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Id         uuid.UUID `json:"id"`
	ProfilePic string    `json:"profliePic"`
}

func CustomUserConvertor(user database.User) CustomUser {
	return CustomUser{
		Username:   user.Username,
		Password:   user.Password,
		Id:         user.ID,
		ProfilePic: user.ProfilePicture.String,
	}
}
