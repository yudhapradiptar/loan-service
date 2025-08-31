package repository

import (
	"context"
	"loan-service/internal/models"

	"gorm.io/gorm"
)

type LoanRepository struct {
	db *gorm.DB
}

// Ensure LoanRepository implements LoanRepositoryInterface
var _ LoanRepositoryInterface = (*LoanRepository)(nil)

func NewLoanRepository(db *gorm.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	return r.db.WithContext(ctx).Begin(), nil
}

func (r *LoanRepository) Commit(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Commit().Error
}

func (r *LoanRepository) Rollback(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Rollback().Error
}

func (r *LoanRepository) CreateLoan(ctx context.Context, loan *models.Loan) error {
	return r.db.WithContext(ctx).Create(loan).Error
}

func (r *LoanRepository) UpdateLoan(ctx context.Context, tx *gorm.DB, loan *models.Loan, fields []string) error {
	return tx.WithContext(ctx).Model(loan).Select(fields).UpdateColumns(loan).Error
}

func (r *LoanRepository) GetLoanByUUID(ctx context.Context, uuid string) (*models.Loan, error) {
	var loan models.Loan
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&loan).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *LoanRepository) GetAllLoans(ctx context.Context) ([]models.Loan, error) {
	var loans []models.Loan
	err := r.db.WithContext(ctx).Find(&loans).Error
	return loans, err
}

func (r *LoanRepository) CreateLoanApproval(ctx context.Context, db *gorm.DB, loanApproval *models.LoanApproval) error {
	return db.WithContext(ctx).Create(loanApproval).Error
}

func (r *LoanRepository) CreateLoanApprovalValidator(ctx context.Context, db *gorm.DB, loanApprovalValidator *models.LoanApprovalValidator) error {
	return db.WithContext(ctx).Create(loanApprovalValidator).Error
}

func (r *LoanRepository) CreateLoanApprovalValidatorProof(ctx context.Context, db *gorm.DB, loanApprovalValidatorProof *models.LoanApprovalValidatorProof) error {
	return db.WithContext(ctx).Create(loanApprovalValidatorProof).Error
}

func (r *LoanRepository) CreateInvestment(ctx context.Context, db *gorm.DB, investment *models.Investment) error {
	return db.WithContext(ctx).Create(investment).Error
}

func (r *LoanRepository) GetInvestmentsByLoanID(ctx context.Context, loanID int) ([]models.Investment, error) {
	var investments []models.Investment
	err := r.db.WithContext(ctx).Where("loan_id = ?", loanID).Find(&investments).Error
	return investments, err
}

func (r *LoanRepository) CreateLoanDisbursement(ctx context.Context, db *gorm.DB, loanDisbursement *models.LoanDisbursement) error {
	return db.WithContext(ctx).Create(loanDisbursement).Error
}
