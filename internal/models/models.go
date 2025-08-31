package models

import (
	"database/sql"
	"loan-service/enums"
	"time"
)

type Loan struct {
	ID               int              `json:"id" gorm:"primaryKey"`
	UUID             string           `json:"uuid" gorm:"not null"`
	BorrowerID       string           `json:"borrower_id" gorm:"not null"`
	PrincipalAmount  float64          `json:"principal_amount" gorm:"not null"`
	InterestRate     float64          `json:"interest_rate" gorm:"not null"`
	ROIRate          float64          `json:"roi_rate" gorm:"not null"`
	InvestmentAmount float64          `json:"investment_amount" gorm:"not null"`
	Status           enums.LoanStatus `json:"status" gorm:"default:1"`
	CreatedAt        time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

type LoanApproval struct {
	ID         int          `json:"id" gorm:"primaryKey"`
	UUID       string       `json:"uuid" gorm:"not null"`
	LoanID     int          `json:"loan_id" gorm:"not null"`
	ApprovedAt sql.NullTime `json:"approved_at"`
	CreatedAt  time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

type LoanApprovalValidator struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	UUID           string    `json:"uuid" gorm:"not null"`
	LoanApprovalID int       `json:"loan_approval_id" gorm:"not null"`
	EmployeeID     string    `json:"employee_id" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type LoanApprovalValidatorProof struct {
	ID                      int       `json:"id" gorm:"primaryKey"`
	UUID                    string    `json:"uuid" gorm:"not null"`
	LoanApprovalValidatorID int       `json:"loan_approval_validator_id" gorm:"not null"`
	ProofURL                string    `json:"proof_url" gorm:"not null"`
	Category                string    `json:"category" gorm:"not null"`
	CreatedAt               time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt               time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Investment struct {
	ID                 int       `json:"id" gorm:"primaryKey"`
	UUID               string    `json:"uuid" gorm:"not null"`
	LoanID             int       `json:"loan_id" gorm:"not null"`
	InvestorID         string    `json:"investor_id" gorm:"not null"`
	Amount             float64   `json:"amount" gorm:"not null"`
	AgreementLetterURL string    `json:"agreement_letter_url" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type LoanDisbursement struct {
	ID                       int       `json:"id" gorm:"primaryKey"`
	UUID                     string    `json:"uuid" gorm:"not null"`
	LoanID                   int       `json:"loan_id" gorm:"not null"`
	FieldOfficerEmployeeID   string    `json:"field_officer_employee_id" gorm:"not null"`
	SignedAgreementLetterURL string    `json:"signed_agreement_letter_url" gorm:"not null"`
	DisbursedAt              time.Time `json:"disbursed_at" gorm:"not null"`
	CreatedAt                time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt                time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
