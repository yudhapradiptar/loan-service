package handlers

import (
	"encoding/json"
	"net/http"

	"loan-service/internal/dto"
	"loan-service/internal/service"

	"github.com/gorilla/mux"
)

type LoanHandler struct {
	loanService *service.LoanService
}

// Ensure LoanHandler implements LoanHandlerInterface
var _ LoanHandlerInterface = (*LoanHandler)(nil)

func NewLoanHandler(loanService *service.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.loanService.CreateLoan(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.APIResponse{
		Message: "Loan created successfully",
	})
}

func (h *LoanHandler) GetAllLoans(w http.ResponseWriter, r *http.Request) {
	loans, err := h.loanService.GetAllLoans(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve loans", http.StatusInternalServerError)
		return
	}

	response := dto.APIResponse{
		Message: "Loans retrieved successfully",
		Data:    loans,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *LoanHandler) GetLoanByUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if uuid == "" {
		http.Error(w, "Missing loan UUID", http.StatusBadRequest)
		return
	}

	loan, err := h.loanService.GetLoanByUUID(r.Context(), uuid)
	if err != nil {
		http.Error(w, "Loan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.APIResponse{
		Message: "Loan retrieved successfully",
		Data:    loan,
	})
}

func (h *LoanHandler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if uuid == "" {
		http.Error(w, "Missing loan UUID", http.StatusBadRequest)
		return
	}

	// TODO: Implement loan approval logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.APIResponse{
		Message: "Loan approval endpoint - not implemented yet",
	})
}

func (h *LoanHandler) InvestLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if uuid == "" {
		http.Error(w, "Missing loan UUID", http.StatusBadRequest)
		return
	}

	// TODO: Implement loan investment logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.APIResponse{
		Message: "Loan investment endpoint - not implemented yet",
	})
}

func (h *LoanHandler) DisburseLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if uuid == "" {
		http.Error(w, "Missing loan UUID", http.StatusBadRequest)
		return
	}

	// TODO: Implement loan disbursement logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.APIResponse{
		Message: "Loan disbursement endpoint - not implemented yet",
	})
}
