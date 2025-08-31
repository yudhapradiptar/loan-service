package service

import (
	"errors"
	"loan-service/internal/models"
	"loan-service/internal/repository"
)

type LoanService struct {
	repo *repository.LoanRepository
}

func NewLoanService(repo *repository.LoanRepository) *LoanService {
	return &LoanService{
		repo: repo,
	}
}

func (s *LoanService) CreateLoan(req *models.CreateLoanRequest) (*models.Loan, error) {
	if req.Amount <= 0 {
		return nil, errors.New("loan amount must be greater than zero")
	}

	if req.InterestRate <= 0 {
		return nil, errors.New("interest rate must be greater than zero")
	}

	if req.Term <= 0 {
		return nil, errors.New("loan term must be greater than zero")
	}

	loan := &models.Loan{
		UserID:       req.UserID,
		Amount:       req.Amount,
		InterestRate: req.InterestRate,
		Term:         req.Term,
		Status:       models.LoanStatusProposed,
	}

	if err := s.repo.CreateLoan(loan); err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *LoanService) GetAllLoans() ([]models.Loan, error) {
	return s.repo.GetAllLoans()
}

func (s *LoanService) GetLoanByID(id int) (*models.Loan, error) {
	return s.repo.GetLoanByID(id)
}

func (s *LoanService) UpdateLoan(id int, req *models.UpdateLoanRequest) (*models.Loan, error) {
	loan, err := s.repo.GetLoanByID(id)
	if err != nil {
		return nil, err
	}

	if req.Amount != nil && *req.Amount <= 0 {
		return nil, errors.New("loan amount must be greater than zero")
	}

	if req.InterestRate != nil && *req.InterestRate <= 0 {
		return nil, errors.New("interest rate must be greater than zero")
	}

	if req.Term != nil && *req.Term <= 0 {
		return nil, errors.New("loan term must be greater than zero")
	}

	if req.Amount != nil {
		loan.Amount = *req.Amount
	}
	if req.InterestRate != nil {
		loan.InterestRate = *req.InterestRate
	}
	if req.Term != nil {
		loan.Term = *req.Term
	}
	if req.Status != nil {
		loan.Status = *req.Status
	}

	if err := s.repo.UpdateLoan(loan); err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *LoanService) DeleteLoan(id int) error {
	return s.repo.DeleteLoan(id)
}
