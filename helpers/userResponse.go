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

type CustomSong struct {
	Name     string    `json:"name"`
	SongLink string    `json:"song_link"`
	Id       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
}

func CustomSongConvertor(song database.Song) CustomSong {
	return CustomSong{
		Name:     song.Name,
		SongLink: song.SongLink,
		Id:       song.ID,
		UserID:   song.UserID,
	}
}
