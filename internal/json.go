package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal payload: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Fatal("Responding with 5XX error: ", message)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	RespondWithJson(w, code, errResponse{
		Error: message,
	})
}