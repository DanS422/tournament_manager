package http

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func WriteResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func WriteValidationError(w http.ResponseWriter, fields map[string]string) {
	WriteResponse(w, http.StatusBadRequest, ErrorResponse{
		Error:  "validation_failed",
		Fields: fields,
	})
}

func WriteError(w http.ResponseWriter, status int, code string) {
	WriteResponse(w, status, ErrorResponse{
		Error: code,
	})
}
