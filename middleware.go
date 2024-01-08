package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/m-rstewart/go-rss/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		handler(w, r, user)
	}
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
