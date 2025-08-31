package repository

import (
	"context"
	"loan-service/internal/models"

	"gorm.io/gorm"
)

type LoanRepositoryInterface interface {
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	Commit(ctx context.Context, db *gorm.DB) error
	Rollback(ctx context.Context, db *gorm.DB) error
	CreateLoan(ctx context.Context, loan *models.Loan) error
	GetLoanByUUID(ctx context.Context, uuid string) (*models.Loan, error)
	GetAllLoans(ctx context.Context) ([]models.Loan, error)
	CreateLoanApproval(ctx context.Context, db *gorm.DB, loanApproval *models.LoanApproval) error
	CreateLoanApprovalValidator(ctx context.Context, db *gorm.DB, loanApprovalValidator *models.LoanApprovalValidator) error
	CreateLoanApprovalValidatorProof(ctx context.Context, db *gorm.DB, loanApprovalValidatorProof *models.LoanApprovalValidatorProof) error
	CreateInvestment(ctx context.Context, db *gorm.DB, investment *models.Investment) error
	UpdateLoan(ctx context.Context, db *gorm.DB, loan *models.Loan, fields []string) error
	GetInvestmentsByLoanID(ctx context.Context, loanID int) ([]models.Investment, error)
	CreateLoanDisbursement(ctx context.Context, db *gorm.DB, loanDisbursement *models.LoanDisbursement) error
}
