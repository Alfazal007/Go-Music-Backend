package controllers

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"
	helper "spotify/helpers"
	"spotify/internal/database"
)

func (apiCfg *ApiConfig) UpdateUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helper.RespondWithError(w, 400, "Issue with finding the user from the database")
		return
	}

	// Parse the form to retrieve file data
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		helper.RespondWithError(w, 400, "Data sent is too big")
		return
	}

	// Retrieve the file from form data
	file, handler, err := r.FormFile("image")
	if err != nil {
		helper.RespondWithError(w, 400, "Error retreiving the file")
		return
	}
	defer file.Close()
	uploadDir := "uploads"
	dstPath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		os.Remove(dstPath)
		helper.RespondWithError(w, 400, "Error saving the file")
		return
	}
	defer dst.Close()

	// Copy the uploaded file's content to the new file
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(dstPath)
		helper.RespondWithError(w, 400, "Error saving the file")
		return
	}
	url, err := apiCfg.cloudinarUploader(r, handler.Filename)
	if err != nil {
		os.Remove(dstPath)
		helper.RespondWithError(w, 400, "Error uploading to cloudinary")
		return
	}
	os.Remove(dstPath)
	updatedUser, err := apiCfg.DB.UpdateProfile(r.Context(), database.UpdateProfileParams{
		ProfilePicture: sql.NullString{String: url, Valid: true},
		ID:             user.ID,
	})
	if err != nil {
		os.Remove(dstPath)
		helper.RespondWithError(w, 400, "Error contacting the database")
		return
	}
	helper.RespondWithJSON(w, 200, helper.CustomUserConvertor(updatedUser))
}
