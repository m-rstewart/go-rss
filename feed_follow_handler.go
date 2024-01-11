package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/m-rstewart/go-rss/internal/database"
)

func (cfg *apiConfig) createFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	type FeedFollowResponse struct {
		ID        uuid.UUID `json:"id"`
		FeedID    uuid.UUID `json:"feed_id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	// Should check to see if feed id exists
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), feedFollowParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create feed follow")
		return
	}

	res := FeedFollowResponse{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID:    feedFollow.UserID,
		FeedID:    feedFollow.FeedID,
	}

	respondWithJSON(w, http.StatusCreated, res)
}

func (cfg *apiConfig) getFeedFollowsHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type GetFeedFollowsResponse struct {
		FeedFollows []database.FeedFollow `json:"feed_follows"`
	}

	feedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get feed follows")
		return
	}

	res := GetFeedFollowsResponse{
		FeedFollows: feedFollows,
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *apiConfig) deleteFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDString := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return
	}

	deleteParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     feedFollowID,
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), deleteParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Feed follow could not be deleted")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
