package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/m-rstewart/go-rss/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_CONN")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error:", err)
	}
	dbQueries := database.New(db)
	apiConfig := &apiConfig{
		DB: dbQueries,
	}

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
	v1Router.Post("/users", apiConfig.createUserHandler)
	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.getCurrentUser))
	v1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.createFeedHandler))
	v1Router.Get("/feeds", apiConfig.getAllFeeds)

	appRouter.Mount("/v1", v1Router)

	fmt.Printf("Starting server on http://localhost%s...\n", server.Addr)
	err = server.ListenAndServe()
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
