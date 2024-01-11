package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/m-rstewart/go-rss/internal/database"
)

type FeedResponse struct {
	Feed       Feed               `json:"feed"`
	FeedFollow FeedFollowResponse `json:"feed_follow"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), feedParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), feedFollowParams)
	if err != nil {
		fmt.Errorf("Error creating feed follow")
	}

	createdFeed := Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}
	createdFeedFollow := FeedFollowResponse{
		ID:        feed.ID,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
	}

	res := FeedResponse{
		Feed:       createdFeed,
		FeedFollow: createdFeedFollow,
	}

	respondWithJSON(w, http.StatusCreated, res)
}

func (cfg *apiConfig) getAllFeeds(w http.ResponseWriter, r *http.Request) {
	type FeedsResponse struct {
		Feeds []database.Feed `json:"feeds"`
	}

	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := FeedsResponse{
		Feeds: feeds,
	}

	respondWithJSON(w, http.StatusOK, res)
}
