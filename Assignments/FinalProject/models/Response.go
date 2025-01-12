package models

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Response interface {
	RespondWithJSON(w http.ResponseWriter, status int, payload interface{})
	RespondWithError(w http.ResponseWriter, code int, message string)
}

func (e ErrorResponse) RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
func (e ErrorResponse) RespondWithError(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	e.RespondWithJSON(w, code, map[string]string{"error": message})
}
