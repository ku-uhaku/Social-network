package helper

import (
	"encoding/json"
	"net/http"
)

// JSONResponse defines the uniform structure for all API outputs
type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSON sends a standardized JSON payload to the client
func WriteJSON(w http.ResponseWriter, status int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := JSONResponse{
		Success: success,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

// Success handles successful operations cleanly
func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	WriteJSON(w, status, true, message, data)
}

// Error handles API errors uniformly
func Error(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, false, message, nil)
}
