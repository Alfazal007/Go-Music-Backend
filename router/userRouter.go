package router

import (
	"net/http"
	"spotify/controllers"

	"github.com/go-chi/chi/v5"
)

func UserRouter(apiCfg *controllers.ApiConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/register", apiCfg.CreateNewUser)
	r.Post("/login", apiCfg.Login)
	r.Post("/refresh-token", apiCfg.RefreshAccessToken)
	r.Get("/current-user", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.GetCurrentUser)).ServeHTTP)
	r.Put("/update-username", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.UpdataUsername)).ServeHTTP)
	r.Put("/update-profile", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.UpdateUserProfilePicture)).ServeHTTP)
	return r
}
