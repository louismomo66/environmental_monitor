package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding data to JSON: %v", err)
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}