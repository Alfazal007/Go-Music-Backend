package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	helper "spotify/helpers"
	"spotify/internal/database"
	"strings"

	"github.com/google/uuid"
)

// register the user
func (apiCfg *ApiConfig) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		helper.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}
	usernameLen := len(strings.TrimSpace(params.Username))
	passwordLen := len(strings.TrimSpace(params.Password))
	if usernameLen < 5 || usernameLen > 30 {
		helper.RespondWithError(w, 400, "Username length should be between 5 and 30")
		return
	}
	if passwordLen < 6 {
		helper.RespondWithError(w, 400, "Password length should be greater than or equal to 6")
		return
	}
	hashedPassword, err := helper.HashPassword(strings.TrimSpace(params.Password))
	if err != nil {
		helper.RespondWithError(w, 400, "There was an issue hashing the password")
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Username: strings.TrimSpace(params.Username),
		ID:       uuid.New(),
		Password: hashedPassword,
	})
	if err != nil {
		helper.RespondWithError(w, 400, fmt.Sprintf("Could not create user %v", err))
		return
	}
	helper.RespondWithJSON(w, 201, helper.CustomUserConvertor(user))
}

func (apiCfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		helper.RespondWithError(w, 400, fmt.Sprintf("There was an error with the request body sent %v", err))
		return
	}
	user, err := apiCfg.DB.GetUserByName(r.Context(), params.Username)
	if err != nil {
		helper.RespondWithJSON(w, 404, "USER NOT FOUND")
		return
	}
	isValid := helper.CheckPasswordHash(params.Password, user.Password)
	if !isValid {
		helper.RespondWithJSON(w, 400, "Incorrect password")
		return
	}
	// send the api key as well in the headers
	jwtToken, err := GenerateJWT(user)
	if err != nil {
		fmt.Println("The error is ", err)
		helper.RespondWithError(w, 400, "Error generating the token")
		return
	}
	cookie := http.Cookie{
		Name:     "access-token",
		Value:    jwtToken,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	// send the api key as well in the headers
	jwtTokenRefresh, err := GenerateRefreshToken(user)
	if err != nil {
		fmt.Println("The error is ", err)
		helper.RespondWithError(w, 400, "Error generating the refresh token")
		return
	}
	_, err = apiCfg.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		RefreshToken: sql.NullString{
			String: jwtTokenRefresh,
			Valid:  true,
		},
		Username: user.Username,
	})
	if err != nil {
		helper.RespondWithError(w, 400, "issue writing token to the database")

		return
	}
	cookie = http.Cookie{
		Name:     "refresh-token",
		Value:    jwtTokenRefresh,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	type AccessToken struct {
		AccessToken  string `json:"access-token"`
		RefreshToken string `json:"refresh-token"`
	}
	helper.RespondWithJSON(w, 200, AccessToken{AccessToken: jwtToken, RefreshToken: jwtTokenRefresh})
}
