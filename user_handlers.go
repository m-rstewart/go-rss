package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m-rstewart/go-rss/internal/database"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
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
		APIKey:    user.ApiKey,
	}

	respondWithJSON(w, http.StatusCreated, res)
}

func (cfg *apiConfig) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := getAPIKeyFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No api key found")
		return
	}

	user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		APIKey:    user.ApiKey,
	}

	respondWithJSON(w, http.StatusOK, res)
}

func getAPIKeyFromRequest(req *http.Request) (string, error) {
	apiKeyHeader := req.Header.Get("Authorization")
	if apiKeyHeader == "" {
		return "", fmt.Errorf("no ApiKey header found")
	}

	keyParts := strings.Split(apiKeyHeader, " ")
	if len(keyParts) != 2 || strings.ToLower(keyParts[0]) != "apikey" {
		return "", fmt.Errorf("invalid ApiKey header format")
	}

	return keyParts[1], nil
}
