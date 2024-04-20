package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rohitdhas/rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	PORT := os.Getenv("PORT")

	if PORT == "" {
		log.Fatal("PORT is not defined!")
	}

	DB_URL := os.Getenv("DB_URL")

	if DB_URL == "" {
		log.Fatal("DB_URL is not defined!")
	}

	conn, err := sql.Open("postgres", DB_URL)

	if err != nil {
		log.Fatal("Error while connecting to postgres!")
	}

	quries := database.New(conn)

	apiCfg := apiConfig{
		DB: quries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	routerV1 := chi.NewRouter()

	routerV1.Get("/healthz", handlerReadiness)
	routerV1.Get("/error", handlerError)
	routerV1.Get("/user", apiCfg.handleGetUser)
	routerV1.Post("/user", apiCfg.handleCreateUser)

	router.Mount("/v1", routerV1)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	log.Printf("Server running on PORT %v", PORT)
	serverErr := server.ListenAndServe()

	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
