package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status":  "ok",
		"service": "loan-service",
	}

	json.NewEncoder(w).Encode(response)
}
