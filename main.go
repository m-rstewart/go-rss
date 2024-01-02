package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	appRouter := chi.NewRouter()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: appRouter,
	}

	corsMiddleware := cors.Handler(cors.Options{})
	appRouter.Use(corsMiddleware)

	v1Router := chi.NewRouter()

	appRouter.Mount("/v1", v1Router)

	fmt.Printf("Starting server on http://localhost%s...\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
