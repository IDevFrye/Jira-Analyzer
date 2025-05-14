package responseutils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, log *slog.Logger, statusCode int, message string, err error) {
	log.Error(message, "error", err)
	writeJSON(w, statusCode, ErrorResponse{Error: message})
}

func WriteErrorSimple(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, ErrorResponse{Error: message})
}

func WriteSuccess(w http.ResponseWriter, log *slog.Logger, statusCode int, payload any) {
	writeJSON(w, statusCode, payload)
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
