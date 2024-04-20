package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/rohitdhas/rssagg/internal"
	"github.com/rohitdhas/rssagg/internal/auth"
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
	}

	internal.RespondWithJson(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)

	if err != nil {
		internal.RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
		return
	}

	type params struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)

	reqParams := params{}
	err = decoder.Decode(&reqParams)

	if err != nil {
		internal.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUserByEmail(r.Context(), reqParams.Email)

	if err != nil {
		internal.RespondWithError(w, 403, fmt.Sprintf("Error fetching user: %v", err))
		return
	}

	if user.ApiKey != apiKey {
		internal.RespondWithError(w, 403, "Invalid api_key")
		return
	}

	internal.RespondWithJson(w, 200, databaseUserToUser(user))
}
