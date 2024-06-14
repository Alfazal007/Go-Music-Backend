package router

import (
	"net/http"
	"spotify/controllers"

	"github.com/go-chi/chi/v5"
)

func SongRouter(apiCfg *controllers.ApiConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/create-song", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.CreateSong)).ServeHTTP)
	r.Get("/get-song/{name}", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.GetSong)).ServeHTTP)
	r.Delete("/delete-song", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.DeleteSong)).ServeHTTP)
	return r
}
