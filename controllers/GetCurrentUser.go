package controllers

import (
	"net/http"
	helper "spotify/helpers"
	"spotify/internal/database"
)

func (apiCfg *ApiConfig) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helper.RespondWithError(w, 400, "Issue with finding the user from the database")
		return
	}
	helper.RespondWithJSON(w, 200, helper.CustomUserConvertor(user))
}
