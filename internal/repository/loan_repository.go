package repository

import (
	"loan-service/internal/models"

	"gorm.io/gorm"
)

type LoanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) CreateLoan(loan *models.Loan) error {
	return r.db.Create(loan).Error
}

func (r *LoanRepository) GetLoanByID(id int) (*models.Loan, error) {
	var loan models.Loan
	err := r.db.First(&loan, id).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *LoanRepository) GetAllLoans() ([]models.Loan, error) {
	var loans []models.Loan
	err := r.db.Find(&loans).Error
	return loans, err
}

func (r *LoanRepository) UpdateLoan(loan *models.Loan) error {
	return r.db.Save(loan).Error
}

func (r *LoanRepository) DeleteLoan(id int) error {
	return r.db.Delete(&models.Loan{}, id).Error
}

func (r *LoanRepository) GetLoansByUserID(userID int) ([]models.Loan, error) {
	var loans []models.Loan
	err := r.db.Where("user_id = ?", userID).Find(&loans).Error
	return loans, err
}
