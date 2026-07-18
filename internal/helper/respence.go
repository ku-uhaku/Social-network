package helper

import (
	"encoding/json"
	"net/http"

	"kuu/internal/requests"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, success bool, message string, data interface{}, errs interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := JSONResponse{
		Success: success,
		Message: message,
		Data:    data,
		Errors:  errs,
	}

	json.NewEncoder(w).Encode(resp)
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	WriteJSON(w, status, true, message, data, nil)
}

func Error(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, false, message, nil, nil)
}

func ValidationErrorResponse(w http.ResponseWriter, status int, errs []requests.ValidationError) {
	WriteJSON(w, status, false, "Validation failed", nil, errs)
}
