package handlers

import (
	"bytes"
	"encoding/json"
	"loan-service/internal/dto"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupTestHandler() *LoanHandler {
	validator := validator.New()
	handler := &LoanHandler{
		loanService: nil,
		validator:   validator,
	}
	return handler
}

func createTestRequest(method, url string, body interface{}) *http.Request {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestLoanHandler_CreateLoan_InvalidJSON(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("POST", "/v1/loans", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateLoan(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestLoanHandler_CreateLoan_ValidationError(t *testing.T) {
	handler := setupTestHandler()

	reqBody := dto.CreateLoanRequest{
		BorrowerID:      "",
		PrincipalAmount: -1000.0,
		InterestRate:    0.05,
		ROIRate:         0.08,
	}

	req := createTestRequest("POST", "/v1/loans", reqBody)
	w := httptest.NewRecorder()

	handler.CreateLoan(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestLoanHandler_GetLoanByUUID_MissingUUID(t *testing.T) {
	handler := setupTestHandler()

	req := createTestRequest("GET", "/v1/loans/", nil)
	w := httptest.NewRecorder()

	handler.GetLoanByUUID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing loan UUID")
}

func TestLoanHandler_ApproveLoan_MissingUUID(t *testing.T) {
	handler := setupTestHandler()

	reqBody := dto.ApproveLoanRequest{
		EmployeeID: "emp1",
		Proofs:     []dto.LoanApprovalValidatorProof{},
		ApprovedAt: time.Now(),
	}

	req := createTestRequest("POST", "/v1/loans/approve", reqBody)
	w := httptest.NewRecorder()

	handler.ApproveLoan(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing loan UUID")
}

func TestLoanHandler_ApproveLoan_InvalidJSON(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("POST", "/v1/loans/test-uuid/approve", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"uuid": "test-uuid"}
	req = mux.SetURLVars(req, vars)

	handler.ApproveLoan(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestLoanHandler_InvestLoan_MissingUUID(t *testing.T) {
	handler := setupTestHandler()

	reqBody := dto.InvestLoanRequest{
		LoanUUID:   "test-uuid",
		InvestorID: "investor123",
		Amount:     500.0,
	}

	req := createTestRequest("POST", "/v1/loans/invest", reqBody)
	w := httptest.NewRecorder()

	handler.InvestLoan(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing loan UUID")
}

func TestLoanHandler_DisburseLoan_MissingUUID(t *testing.T) {
	handler := setupTestHandler()

	reqBody := dto.CreateLoanDisbursementRequest{
		LoanUUID:                 "test-uuid",
		EmployeeID:               "emp123",
		SignedAgreementLetterURL: "https://example.com/signed.pdf",
		DisbursedAt:              time.Now(),
	}

	req := createTestRequest("POST", "/v1/loans/disburse", reqBody)
	w := httptest.NewRecorder()

	handler.DisburseLoan(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing loan UUID")
}

func TestLoanHandler_DisburseLoan_ValidationError(t *testing.T) {
	handler := setupTestHandler()

	reqBody := dto.CreateLoanDisbursementRequest{
		LoanUUID:                 "test-uuid",
		EmployeeID:               "",
		SignedAgreementLetterURL: "",
		DisbursedAt:              time.Now(),
	}

	req := createTestRequest("POST", "/v1/loans/test-uuid/disburse", reqBody)
	w := httptest.NewRecorder()

	vars := map[string]string{"uuid": "test-uuid"}
	req = mux.SetURLVars(req, vars)

	handler.DisburseLoan(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}
