package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, errResp ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResp.Status)
	json.NewEncoder(w).Encode(errResp)
}
