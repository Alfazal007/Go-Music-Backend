package controllers

import (
	"spotify/internal/database"

	"github.com/cloudinary/cloudinary-go/v2"
)

type ApiConfig struct {
	DB  *database.Queries
	Cld *cloudinary.Cloudinary
}
