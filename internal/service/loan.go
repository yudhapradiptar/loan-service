package service

import (
	"context"
	"database/sql"
	"errors"
	"loan-service/enums"
	"loan-service/internal/client"
	"loan-service/internal/dto"
	"loan-service/internal/models"
	"loan-service/internal/repository"
)

type LoanService struct {
	repo               repository.LoanRepositoryInterface
	notificationClient client.NotificationClientInterface
}

// Ensure LoanService implements LoanServiceInterface
var _ LoanServiceInterface = (*LoanService)(nil)

func NewLoanService(repo repository.LoanRepositoryInterface, notificationClient client.NotificationClientInterface) *LoanService {
	return &LoanService{
		repo:               repo,
		notificationClient: notificationClient,
	}
}

func (s *LoanService) CreateLoan(ctx context.Context, req *dto.CreateLoanRequest) error {
	loan := &models.Loan{
		BorrowerID:      req.BorrowerID,
		PrincipalAmount: req.PrincipalAmount,
		InterestRate:    req.InterestRate,
		ROIRate:         req.ROIRate,
		Status:          enums.LoanStatusProposed,
	}

	if err := s.repo.CreateLoan(ctx, loan); err != nil {
		return err
	}

	return nil
}

func (s *LoanService) GetAllLoans(ctx context.Context) (response []dto.GetLoansResponseItem, err error) {
	loans, err := s.repo.GetAllLoans(ctx)
	if err != nil {
		return nil, err
	}

	for _, loan := range loans {
		response = append(response, dto.GetLoansResponseItem{
			UUID:            loan.UUID,
			BorrowerID:      loan.BorrowerID,
			PrincipalAmount: loan.PrincipalAmount,
			InterestRate:    loan.InterestRate,
			ROIRate:         loan.ROIRate,
			Status:          loan.Status,
		})
	}

	return response, nil
}

func (s *LoanService) GetLoanByUUID(ctx context.Context, uuid string) (dto.GetLoansResponseItem, error) {
	loan, err := s.repo.GetLoanByUUID(ctx, uuid)
	if err != nil {
		return dto.GetLoansResponseItem{}, err
	}
	return dto.GetLoansResponseItem{
		UUID:            loan.UUID,
		BorrowerID:      loan.BorrowerID,
		PrincipalAmount: loan.PrincipalAmount,
		InterestRate:    loan.InterestRate,
		ROIRate:         loan.ROIRate,
		Status:          loan.Status,
	}, nil
}

func (s *LoanService) ApproveLoanWithValidators(ctx context.Context, req dto.ApproveLoanRequest) error {
	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			s.repo.Rollback(ctx, tx)
		}
	}()

	loan, err := s.repo.GetLoanByUUID(ctx, req.LoanUUID)
	if err != nil {
		return err
	}

	loanApproval := &models.LoanApproval{
		LoanID:     loan.ID,
		ApprovedAt: sql.NullTime{Time: req.ApprovedAt, Valid: true},
	}

	if err := s.repo.CreateLoanApproval(ctx, tx, loanApproval); err != nil {
		return err
	}

	loanApprovalValidator := &models.LoanApprovalValidator{
		LoanApprovalID: loanApproval.ID,
		EmployeeID:     req.EmployeeID,
	}

	if err := s.repo.CreateLoanApprovalValidator(ctx, tx, loanApprovalValidator); err != nil {
		return err
	}

	for _, proof := range req.Proofs {
		if err := s.repo.CreateLoanApprovalValidatorProof(ctx, tx, &models.LoanApprovalValidatorProof{
			LoanApprovalValidatorID: loanApprovalValidator.ID,
			ProofURL:                proof.ProofURL,
			Category:                proof.Category,
		}); err != nil {
			return err
		}
	}

	loan.Status = enums.LoanStatusApproved
	if err := s.repo.UpdateLoan(ctx, tx, loan, []string{"status"}); err != nil {
		return err
	}

	return s.repo.Commit(ctx, tx)
}

func (s *LoanService) InvestLoan(ctx context.Context, req dto.InvestLoanRequest) error {
	loan, err := s.repo.GetLoanByUUID(ctx, req.LoanUUID)
	if err != nil {
		return err
	}

	if loan.Status != enums.LoanStatusApproved {
		return errors.New("loan is not approved")
	}

	if loan.InvestmentAmount+req.Amount > loan.PrincipalAmount {
		return errors.New("loan investment amount is greater than principal amount")
	}

	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			s.repo.Rollback(ctx, tx)
		}
	}()

	loan.InvestmentAmount += req.Amount
	if err := s.repo.UpdateLoan(ctx, tx, loan, []string{"investment_amount"}); err != nil {
		return err
	}

	investment := &models.Investment{
		LoanID:     loan.ID,
		InvestorID: req.InvestorID,
		Amount:     req.Amount,
	}

	agreementLetterURL, err := generateAgreementLetterURL(ctx, loan, investment)
	if err != nil {
		return err
	}

	investment.AgreementLetterURL = agreementLetterURL

	if err := s.repo.CreateInvestment(ctx, tx, investment); err != nil {
		return err
	}

	if loan.InvestmentAmount == loan.PrincipalAmount {
		loan.Status = enums.LoanStatusInvested
		if err := s.repo.UpdateLoan(ctx, tx, loan, []string{"status"}); err != nil {
			return err
		}
	}

	err = s.repo.Commit(ctx, tx)
	if err != nil {
		return err
	}

	// get all investment records
	investments, err := s.repo.GetInvestmentsByLoanID(ctx, loan.ID)
	if err != nil {
		return err
	}

	for _, investment := range investments {
		// send agreement letter attached to email to investor
		err = s.notificationClient.SendEmail(ctx, client.SendEmailRequest{
			To:      investment.InvestorID, // notification service will get the email from the investor id
			Subject: "Loan Agreement Letter",
			Body:    "Please find the agreement letter attached to this email.",
			Attachments: []client.Attachment{
				{
					Filename: "agreement_letter.pdf",
					Content:  investment.AgreementLetterURL,
					Type:     "application/pdf",
				},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func generateAgreementLetterURL(_ context.Context, _ *models.Loan, _ *models.Investment) (string, error) {
	// generate agreement letter
	return "", nil
}

func (s *LoanService) CreateLoanDisbursement(ctx context.Context, req dto.CreateLoanDisbursementRequest) error {
	loan, err := s.repo.GetLoanByUUID(ctx, req.LoanUUID)
	if err != nil {
		return err
	}

	if loan.Status != enums.LoanStatusInvested {
		return errors.New("loan is not invested")
	}

	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			s.repo.Rollback(ctx, tx)
		}
	}()

	loanDisbursement := &models.LoanDisbursement{
		LoanID:                   loan.ID,
		FieldOfficerEmployeeID:   req.EmployeeID,
		SignedAgreementLetterURL: req.SignedAgreementLetterURL,
		DisbursedAt:              req.DisbursedAt,
	}

	if err := s.repo.CreateLoanDisbursement(ctx, tx, loanDisbursement); err != nil {
		return err
	}

	loan.Status = enums.LoanStatusDisbursed
	if err := s.repo.UpdateLoan(ctx, tx, loan, []string{"status"}); err != nil {
		return err
	}

	err = s.repo.Commit(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}
