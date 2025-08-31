package models

import (
	"time"
)

type APIResponse struct {
	Object  string      `json:"object"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ListResponse struct {
	Object  string        `json:"object"`
	Data    []interface{} `json:"data"`
	HasMore bool          `json:"has_more"`
	URL     string        `json:"url"`
}

type LoanResponse struct {
	ID           int       `json:"id"`
	Object       string    `json:"object"`
	UserID       int       `json:"user_id"`
	Amount       float64   `json:"amount"`
	InterestRate float64   `json:"interest_rate"`
	Term         int       `json:"term"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewLoanResponse(loan *Loan) *LoanResponse {
	return &LoanResponse{
		ID:           loan.ID,
		Object:       "loan",
		UserID:       loan.UserID,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		Term:         loan.Term,
		Status:       loan.Status.String(),
		CreatedAt:    loan.CreatedAt,
		UpdatedAt:    loan.UpdatedAt,
	}
}
