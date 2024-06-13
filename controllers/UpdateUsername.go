package controllers

import (
	"encoding/json"
	"net/http"
	helper "spotify/helpers"
	"spotify/internal/database"
)

func (apiCfg *ApiConfig) UpdataUsername(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helper.RespondWithError(w, 400, "Issue with finding the user from the database")
		return
	}
	type parameters struct {
		Username string `json:"username"`
	}
	var params parameters
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)
	user, err := apiCfg.DB.UpdateUsername(r.Context(), database.UpdateUsernameParams{
		Username: params.Username,
		ID:       user.ID,
	})
	if err != nil {
		helper.RespondWithError(w, 400, "Could not update the username check request body")
		return
	}
	helper.RespondWithJSON(w, 200, helper.CustomUserConvertor(user))
}
