package handlers

import (
	"net/http"

	"github.com/rohitdhas/rssagg/internal"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	type healthResponse struct {
		Status string `json:"status"`
	}
	internal.RespondWithJson(w, 200, healthResponse{
		Status: "Running ✅",
	})
}

func HandlerError(w http.ResponseWriter, r *http.Request) {
	internal.RespondWithError(w, 400, "Something went wrong!")
}