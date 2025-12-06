package api

import (
	"encoding/json"
	"io"
	"net/http"

	"tasks-api/internal/models"
)

func ParseBody[T models.Validator](r *http.Request) (*T, bool) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, false
	}

	var obj T
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		return nil, false
	}

	if !obj.Validate() {
		return nil, false
	}

	return &obj, true
}

func SerializeResponse(w http.ResponseWriter, result any) {
	// // short (sufficient) syntax
	// if err := json.NewEncoder(w).Encode(result); err != nil {
	// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// 	return
	// }

	// explicit syntax
	serializedData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Write(serializedData)
}
