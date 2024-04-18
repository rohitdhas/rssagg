package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/rohitdhas/rssagg/handlers"
)

func main() {
	godotenv.Load()
	PORT := os.Getenv("PORT")

	if PORT == "" {
		log.Fatal("PORT is not defined!")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	routerV1 := chi.NewRouter()

	routerV1.Get("/healthz", handlers.HandlerReadiness)
	routerV1.Get("/error", handlers.HandlerError)
	
	router.Mount("/v1", routerV1)

	server := &http.Server{
		Handler: router,
		Addr: ":" + PORT,
	}

	log.Printf("Server running on PORT %v", PORT)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}