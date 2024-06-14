package controllers

import (
	"encoding/json"
	"net/http"
	helper "spotify/helpers"
	"spotify/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) DeleteSong(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helper.RespondWithError(w, 400, "Issue with finding the user from the database")
		return
	}
	type Parameters struct {
		ID string `json:"id"`
	}
	var params Parameters
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)
	song_id, err := uuid.Parse(params.ID)
	if err != nil {
		helper.RespondWithError(w, 400, "Invalid input data")
		return
	}
	_, err = apiCfg.DB.DeleteSong(r.Context(), database.DeleteSongParams{
		ID:     song_id,
		UserID: user.ID,
	})
	if err != nil {
		helper.RespondWithError(w, 400, "Error deleting the song from the database")
		return
	}
	helper.RespondWithJSON(w, 200, "deleted song successfully")
}
