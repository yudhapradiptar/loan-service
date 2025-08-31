package dto

import (
	"loan-service/enums"
	"time"
)

type CreateLoanRequest struct {
	BorrowerID      string  `json:"user_id" validate:"required"`
	PrincipalAmount float64 `json:"principal_amount" validate:"required,gt=0"`
	InterestRate    float64 `json:"interest_rate" validate:"required,gt=0"`
	ROIRate         float64 `json:"roi_rate" validate:"required,gt=0"`
}

type GetLoansResponseItem struct {
	UUID            string           `json:"uuid"`
	BorrowerID      string           `json:"borrower_id"`
	PrincipalAmount float64          `json:"principal_amount"`
	InterestRate    float64          `json:"interest_rate"`
	ROIRate         float64          `json:"roi_rate"`
	Status          enums.LoanStatus `json:"status"`
}

type ApproveLoanRequest struct {
	LoanUUID   string                       `json:"-"`
	EmployeeID string                       `json:"employee_id" validate:"required"`
	Proofs     []LoanApprovalValidatorProof `json:"proofs" validate:"required"`
	ApprovedAt time.Time                    `json:"approved_at" validate:"required"`
}

type LoanApprovalValidatorProof struct {
	ProofURL string `json:"proof_url" validate:"required"`
	Category string `json:"category" validate:"required"`
}

type InvestLoanRequest struct {
	LoanUUID   string  `json:"loan_uuid" validate:"required"`
	InvestorID string  `json:"investor_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

type CreateLoanDisbursementRequest struct {
	LoanUUID                 string    `json:"loan_uuid" validate:"required"`
	EmployeeID               string    `json:"employee_id" validate:"required"`
	SignedAgreementLetterURL string    `json:"signed_agreement_letter_url" validate:"required"`
	DisbursedAt              time.Time `json:"disbursed_at" validate:"required"`
}
