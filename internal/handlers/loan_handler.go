package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"loan-service/internal/models"
	"loan-service/internal/service"

	"github.com/gorilla/mux"
)

type LoanHandler struct {
	loanService *service.LoanService
}

func NewLoanHandler(loanService *service.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	loan, err := h.loanService.CreateLoan(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.NewLoanResponse(loan))
}

func (h *LoanHandler) GetAllLoans(w http.ResponseWriter, r *http.Request) {
	loans, err := h.loanService.GetAllLoans()
	if err != nil {
		http.Error(w, "Failed to retrieve loans", http.StatusInternalServerError)
		return
	}

	loanResponses := make([]*models.LoanResponse, len(loans))
	for i, loan := range loans {
		loanResponses[i] = models.NewLoanResponse(&loan)
	}

	response := models.ListResponse{
		Object:  "list",
		Data:    make([]interface{}, len(loanResponses)),
		HasMore: false,
		URL:     "/v1/loans",
	}

	for i, loanResp := range loanResponses {
		response.Data[i] = loanResp
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *LoanHandler) GetLoanByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		http.Error(w, "Missing loan ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	loan, err := h.loanService.GetLoanByID(id)
	if err != nil {
		http.Error(w, "Loan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.NewLoanResponse(loan))
}

func (h *LoanHandler) UpdateLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		http.Error(w, "Missing loan ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	loan, err := h.loanService.UpdateLoan(id, &req)
	if err != nil {
		http.Error(w, "Loan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.NewLoanResponse(loan))
}

func (h *LoanHandler) DeleteLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		http.Error(w, "Missing loan ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	err = h.loanService.DeleteLoan(id)
	if err != nil {
		http.Error(w, "Loan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.APIResponse{
		Object:  "deleted",
		Message: "Loan deleted successfully",
	})
}
