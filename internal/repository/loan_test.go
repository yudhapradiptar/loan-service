package repository

import (
	"context"
	"loan-service/internal/models"
	"testing"
	"time"

	"database/sql"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Loan{}, &models.LoanApproval{}, &models.LoanApprovalValidator{},
		&models.LoanApprovalValidatorProof{}, &models.Investment{}, &models.LoanDisbursement{})
	assert.NoError(t, err)

	return db
}

func TestLoanRepository_CreateLoan(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()
	loan := &models.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 1000.0,
		InterestRate:    0.05,
		ROIRate:         0.08,
	}

	err := repo.CreateLoan(ctx, loan)
	assert.NoError(t, err)
	assert.NotEmpty(t, loan.UUID)
	assert.NotZero(t, loan.ID)
}

func TestLoanRepository_GetLoanByUUID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()

	loan := &models.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 1000.0,
		InterestRate:    0.05,
		ROIRate:         0.08,
	}

	err := repo.CreateLoan(ctx, loan)
	assert.NoError(t, err)

	retrievedLoan, err := repo.GetLoanByUUID(ctx, loan.UUID)
	assert.NoError(t, err)
	assert.Equal(t, loan.UUID, retrievedLoan.UUID)
	assert.Equal(t, loan.BorrowerID, retrievedLoan.BorrowerID)
}

func TestLoanRepository_GetLoanByUUID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()
	loan, err := repo.GetLoanByUUID(ctx, "non-existent-uuid")
	assert.Error(t, err)
	assert.Nil(t, loan)
}

func TestLoanRepository_GetAllLoans(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()

	loan1 := &models.Loan{
		BorrowerID:      "user1",
		PrincipalAmount: 1000.0,
		InterestRate:    0.05,
		ROIRate:         0.08,
	}

	loan2 := &models.Loan{
		BorrowerID:      "user2",
		PrincipalAmount: 2000.0,
		InterestRate:    0.06,
		ROIRate:         0.09,
	}

	err := repo.CreateLoan(ctx, loan1)
	assert.NoError(t, err)
	err = repo.CreateLoan(ctx, loan2)
	assert.NoError(t, err)

	loans, err := repo.GetAllLoans(ctx)
	assert.NoError(t, err)
	assert.Len(t, loans, 2)
}

func TestLoanRepository_UpdateLoan(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()
	loan := &models.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 1000.0,
		InterestRate:    0.05,
		ROIRate:         0.08,
	}

	err := repo.CreateLoan(ctx, loan)
	assert.NoError(t, err)

	tx, err := repo.BeginTransaction(ctx)
	assert.NoError(t, err)

	loan.PrincipalAmount = 1500.0
	err = repo.UpdateLoan(ctx, tx, loan, []string{"principal_amount"})
	assert.NoError(t, err)

	err = repo.Commit(ctx, tx)
	assert.NoError(t, err)

	retrievedLoan, err := repo.GetLoanByUUID(ctx, loan.UUID)
	assert.NoError(t, err)
	assert.Equal(t, 1500.0, retrievedLoan.PrincipalAmount)
}

func TestLoanRepository_CreateLoanApproval(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()
	tx, err := repo.BeginTransaction(ctx)
	assert.NoError(t, err)

	loanApproval := &models.LoanApproval{
		LoanID:     1,
		ApprovedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	err = repo.CreateLoanApproval(ctx, tx, loanApproval)
	assert.NoError(t, err)
	assert.NotEmpty(t, loanApproval.UUID)
	assert.NotZero(t, loanApproval.ID)

	err = repo.Commit(ctx, tx)
	assert.NoError(t, err)
}

func TestLoanRepository_CreateLoanApprovalValidator(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()
	tx, err := repo.BeginTransaction(ctx)
	assert.NoError(t, err)

	validator := &models.LoanApprovalValidator{
		LoanApprovalID: 1,
		EmployeeID:     "emp123",
	}

	err = repo.CreateLoanApprovalValidator(ctx, tx, validator)
	assert.NoError(t, err)
	assert.NotEmpty(t, validator.UUID)
	assert.NotZero(t, validator.ID)

	err = repo.Commit(ctx, tx)
	assert.NoError(t, err)
}

func TestLoanRepository_CreateLoanApprovalValidatorProof(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()
	tx, err := repo.BeginTransaction(ctx)
	assert.NoError(t, err)

	proof := &models.LoanApprovalValidatorProof{
		LoanApprovalValidatorID: 1,
		ProofURL:                "https://example.com/proof.pdf",
		Category:                "identity",
	}

	err = repo.CreateLoanApprovalValidatorProof(ctx, tx, proof)
	assert.NoError(t, err)
	assert.NotEmpty(t, proof.UUID)
	assert.NotZero(t, proof.ID)

	err = repo.Commit(ctx, tx)
	assert.NoError(t, err)
}

func TestLoanRepository_CreateInvestment(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()
	tx, err := repo.BeginTransaction(ctx)
	assert.NoError(t, err)

	investment := &models.Investment{
		LoanID:     1,
		InvestorID: "investor123",
		Amount:     500.0,
	}

	err = repo.CreateInvestment(ctx, tx, investment)
	assert.NoError(t, err)
	assert.NotEmpty(t, investment.UUID)
	assert.NotZero(t, investment.ID)

	err = repo.Commit(ctx, tx)
	assert.NoError(t, err)
}

func TestLoanRepository_GetInvestmentsByLoanID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()
	tx, err := repo.BeginTransaction(ctx)
	assert.NoError(t, err)

	investment1 := &models.Investment{
		LoanID:     1,
		InvestorID: "investor1",
		Amount:     500.0,
	}

	investment2 := &models.Investment{
		LoanID:     1,
		InvestorID: "investor2",
		Amount:     300.0,
	}

	err = repo.CreateInvestment(ctx, tx, investment1)
	assert.NoError(t, err)
	err = repo.CreateInvestment(ctx, tx, investment2)
	assert.NoError(t, err)

	err = repo.Commit(ctx, tx)
	assert.NoError(t, err)

	investments, err := repo.GetInvestmentsByLoanID(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, investments, 2)
}

func TestLoanRepository_CreateLoanDisbursement(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()
	tx, err := repo.BeginTransaction(ctx)
	assert.NoError(t, err)

	disbursement := &models.LoanDisbursement{
		LoanID:                   1,
		FieldOfficerEmployeeID:   "emp123",
		SignedAgreementLetterURL: "https://example.com/signed.pdf",
		DisbursedAt:              time.Now(),
	}

	err = repo.CreateLoanDisbursement(ctx, tx, disbursement)
	assert.NoError(t, err)
	assert.NotEmpty(t, disbursement.UUID)
	assert.NotZero(t, disbursement.ID)

	err = repo.Commit(ctx, tx)
	assert.NoError(t, err)
}

func TestLoanRepository_TransactionOperations(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLoanRepository(db)

	ctx := context.Background()

	tx, err := repo.BeginTransaction(ctx)
	assert.NoError(t, err)

	err = repo.Commit(ctx, tx)
	assert.NoError(t, err)

	tx2, err := repo.BeginTransaction(ctx)
	assert.NoError(t, err)

	err = repo.Rollback(ctx, tx2)
	assert.NoError(t, err)
}
