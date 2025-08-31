package handlers

import "net/http"

type LoanHandlerInterface interface {
	CreateLoan(w http.ResponseWriter, r *http.Request)
	GetAllLoans(w http.ResponseWriter, r *http.Request)
	GetLoanByUUID(w http.ResponseWriter, r *http.Request)
	ApproveLoan(w http.ResponseWriter, r *http.Request)
	InvestLoan(w http.ResponseWriter, r *http.Request)
	DisburseLoan(w http.ResponseWriter, r *http.Request)
}

type HealthHandlerInterface interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
}
