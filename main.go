package main

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	appRouter := chi.NewRouter()

	corsMiddleware := cors.Handler(cors.Options{})
	appRouter.Use(corsMiddleware)

	v1Router := chi.NewRouter()

	appRouter.Mount("/v1", v1Router)
}
