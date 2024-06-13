package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	helper "spotify/helpers"
	"spotify/internal/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		RefreshToken string `json:"refresh-token"`
	}
	var params parameter
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)
	// Verify the JWT token
	jwtSecret := os.Getenv("SECRET_KEY_REFRESH")
	if jwtSecret == "" {
		helper.RespondWithError(w, 400, "Server error")
		return
	}
	jwtToken := params.RefreshToken
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		helper.RespondWithError(w, 401, fmt.Sprintf("Invalid token here %v", err))
		return
	}
	if !token.Valid {
		helper.RespondWithError(w, 401, fmt.Sprintf("Invalid token %v", err))
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helper.RespondWithError(w, 400, "Invalid claims login again")
		return
	}

	username := claims["username"].(string)
	id := claims["user_id"].(string)

	user, err := apiCfg.DB.GetUserByName(r.Context(), username)
	if err != nil {
		helper.RespondWithError(w, 400, "Some manpulation done with the token")
		return
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		helper.RespondWithError(w, 400, "Some manpulation done with the token")
		return
	}
	if idUUID != user.ID {
		helper.RespondWithError(w, 400, "Some manipulations done with the token try again")
		return
	}
	if params.RefreshToken != user.RefreshToken.String {
		helper.RespondWithError(w, 400, "Invalid refresh token")
		return
	}
	accessTokenNew, err := GenerateJWT(user)
	if err != nil {
		helper.RespondWithError(w, 400, "Error generating access token")
		return
	}
	refreshToken, err := GenerateRefreshToken(user)
	if err != nil {
		helper.RespondWithError(w, 400, "Error generating refresh token")
	}
	user, err = apiCfg.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		RefreshToken: sql.NullString{String: refreshToken, Valid: true},
		Username:     username,
	})
	if err != nil {
		helper.RespondWithError(w, 400, "Could not refresh the tokens")
		return
	}
	cookie := http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:     "access-token",
		Value:    accessTokenNew,
		HttpOnly: true,
	}
	type AccessToken struct {
		AccessToken  string `json:"access-token"`
		RefreshToken string `json:"refresh-token"`
	}
	helper.RespondWithJSON(w, 200, AccessToken{AccessToken: accessTokenNew, RefreshToken: refreshToken})
}
