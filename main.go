package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"spotify/controllers"
	"spotify/internal/database"
	"spotify/router"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		log.Fatal("There was an issue getting port number from the env variables")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB url is not found in env variables")
	}
	conn, err := sql.Open("postgres", dbURL)
	apiCfg := controllers.ApiConfig{DB: database.New(conn)}
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	cloudApiKey := os.Getenv("CLOUDINARY_API_KEY")
	cloudApiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	if cloudName == "" || cloudApiKey == "" || cloudApiSecret == "" {
		log.Fatal("There was an issue getting env variables of cloudinary")
	}

	// Configure your Cloudinary credentials
	cld, err := cloudinary.NewFromParams(cloudName, cloudApiKey, cloudApiSecret)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary, %v", err)
	}
	apiCfg.Cld = cld

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	userRouter := router.UserRouter(&apiCfg)
	r.Mount("/user", userRouter)
	srv := &http.Server{
		Addr:    ":" + portNumber,
		Handler: r,
	}

	fmt.Println("Server listening at port:", portNumber)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("There was an error with the server", err)
	}
}
