package controllers

import (
	"net/http"
	helper "spotify/helpers"
	"spotify/internal/database"

	"github.com/go-chi/chi/v5"
)

func (apiCfg *ApiConfig) GetSong(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user").(database.User)
	if !ok {
		helper.RespondWithError(w, 400, "Issue with finding the user from the database")
		return
	}
	name := chi.URLParam(r, "name")

	song, err := apiCfg.DB.GetSongByName(r.Context(), name)
	if err != nil {
		helper.RespondWithError(w, 400, "Error talking to database or song not found")
		return
	}
	helper.RespondWithJSON(w, 200, helper.CustomSongConvertor(song))
}
