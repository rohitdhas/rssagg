package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/rohitdhas/rssagg/internal"
	"github.com/rohitdhas/rssagg/internal/database"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type healthResponse struct {
		Status string `json:"status"`
	}
	internal.RespondWithJson(w, 200, healthResponse{
		Status: "Running ✅",
	})
}

func handlerError(w http.ResponseWriter, r *http.Request) {
	internal.RespondWithError(w, 400, "Something went wrong!")
}

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)

	reqParams := params{}
	err := decoder.Decode(&reqParams)

	if err != nil {
		internal.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:    uuid.New(),
		Name:  reqParams.Name,
		Email: reqParams.Email,
	})

	if err != nil {
		internal.RespondWithError(w, 500, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	internal.RespondWithJson(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	internal.RespondWithJson(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	reqParams := params{}
	err := decoder.Decode(&reqParams)

	if err != nil {
		internal.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   reqParams.Name,
		Url:    reqParams.Url,
		UserID: user.ID,
	})

	if err != nil {
		internal.RespondWithError(w, 500, fmt.Sprintf("Error creating feed: %v", err))
		return
	}

	internal.RespondWithJson(w, 200, databaseFeedToFeed(feed))
}
