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
	v1Router.Get("/readiness", readinessHandler)
	v1Router.Get("/err", errHandler)

	appRouter.Mount("/v1", v1Router)

	fmt.Printf("Starting server on http://localhost%s...\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	type ReadinessResponse struct {
		Status string `json:"status"`
	}

	respondWithJSON(w, http.StatusOK, ReadinessResponse{Status: "ok"})
}

func errHandler(w http.ResponseWriter, r *http.Request) {
	type ErrResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, http.StatusInternalServerError, ErrResponse{Error: "Internal Server Error"})
}
