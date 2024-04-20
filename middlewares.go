package main

import (
	"fmt"
	"net/http"

	"github.com/rohitdhas/rssagg/internal"
	"github.com/rohitdhas/rssagg/internal/auth"
	"github.com/rohitdhas/rssagg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)

		if err != nil {
			internal.RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			internal.RespondWithError(w, 403, fmt.Sprintf("Error while fetching user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
