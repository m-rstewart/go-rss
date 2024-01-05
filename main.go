package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
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

	appRouter.Mount("/v1", v1Router)

	fmt.Printf("Starting server on http://localhost%s...\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"user"`
	}

	type UserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	}
	user, err := cfg.DB.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
	}

	respondWithJSON(w, http.StatusCreated, res)
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
