package service

import (
	"context"
	"loan-service/internal/dto"
)

type LoanServiceInterface interface {
	CreateLoan(ctx context.Context, req *dto.CreateLoanRequest) error
	GetAllLoans(ctx context.Context) ([]dto.GetLoansResponseItem, error)
	GetLoanByUUID(ctx context.Context, uuid string) (dto.GetLoansResponseItem, error)
	ApproveLoanWithValidators(ctx context.Context, req dto.ApproveLoanRequest) error
	InvestLoan(ctx context.Context, req dto.InvestLoanRequest) error
	CreateLoanDisbursement(ctx context.Context, req dto.CreateLoanDisbursementRequest) error
}
