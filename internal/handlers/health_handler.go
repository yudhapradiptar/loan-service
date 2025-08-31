package handlers

import (
	"encoding/json"
	"loan-service/internal/dto"
	"net/http"
)

type HealthHandler struct{}

// Ensure HealthHandler implements HealthHandlerInterface
var _ HealthHandlerInterface = (*HealthHandler)(nil)

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := dto.APIResponse{
		Message: "Service is healthy",
		Data: map[string]string{
			"status":  "ok",
			"service": "loan-service",
		},
	}

	json.NewEncoder(w).Encode(response)
}
