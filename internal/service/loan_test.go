package service

import (
	"context"
	"errors"
	"loan-service/enums"
	"loan-service/internal/client"
	"loan-service/internal/dto"
	"loan-service/internal/models"
	"loan-service/internal/repository"
	"loan-service/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestLoanService_CreateLoan(t *testing.T) {
	type fields struct {
		repo               repository.LoanRepositoryInterface
		notificationClient client.NotificationClientInterface
	}
	type args struct {
		ctx context.Context
		req *dto.CreateLoanRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("CreateLoan", context.Background(), &models.Loan{
						BorrowerID:      "123",
						PrincipalAmount: 100000000,
						InterestRate:    5,
						ROIRate:         3,
						Status:          enums.LoanStatusProposed,
					}).Return(nil)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &dto.CreateLoanRequest{
					BorrowerID:      "123",
					PrincipalAmount: 100000000,
					InterestRate:    5,
					ROIRate:         3,
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("CreateLoan", context.Background(), &models.Loan{
						BorrowerID:      "123",
						PrincipalAmount: 100000000,
						InterestRate:    5,
						ROIRate:         3,
						Status:          enums.LoanStatusProposed,
					}).Return(errors.New("error"))
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &dto.CreateLoanRequest{
					BorrowerID:      "123",
					PrincipalAmount: 100000000,
					InterestRate:    5,
					ROIRate:         3,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LoanService{
				repo:               tt.fields.repo,
				notificationClient: tt.fields.notificationClient,
			}
			if err := s.CreateLoan(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("LoanService.CreateLoan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoanService_ApproveLoanWithValidators(t *testing.T) {
	type fields struct {
		repo               repository.LoanRepositoryInterface
		notificationClient client.NotificationClientInterface
	}
	type args struct {
		ctx context.Context
		req dto.ApproveLoanRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("BeginTransaction", context.Background()).Return(&gorm.DB{}, nil)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(&models.Loan{
						ID:     1,
						UUID:   "loan-uuid-123",
						Status: enums.LoanStatusProposed,
					}, nil)
					m.On("CreateLoanApproval", context.Background(), mock.Anything, mock.MatchedBy(func(approval *models.LoanApproval) bool {
						return approval.LoanID == 1
					})).Return(nil)
					m.On("CreateLoanApprovalValidator", context.Background(), mock.Anything, mock.MatchedBy(func(validator *models.LoanApprovalValidator) bool {
						return validator.EmployeeID == "emp123"
					})).Return(nil)
					m.On("CreateLoanApprovalValidatorProof", context.Background(), mock.Anything, mock.Anything).Return(nil)
					m.On("UpdateLoan", context.Background(), mock.Anything, mock.Anything, []string{"status"}).Return(nil)
					m.On("Commit", context.Background(), mock.Anything).Return(nil)
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.ApproveLoanRequest{
					LoanUUID:   "loan-uuid-123",
					EmployeeID: "emp123",
					ApprovedAt: time.Now(),
					Proofs: []dto.LoanApprovalValidatorProof{
						{
							ProofURL: "https://example.com/proof1.pdf",
							Category: "identity",
						},
						{
							ProofURL: "https://example.com/proof2.pdf",
							Category: "income",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error - loan not found",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("BeginTransaction", context.Background()).Return(&gorm.DB{}, nil)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(nil, errors.New("loan not found"))
					m.On("Rollback", context.Background(), mock.Anything).Return(nil)
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.ApproveLoanRequest{
					LoanUUID:   "loan-uuid-123",
					EmployeeID: "emp123",
					ApprovedAt: time.Now(),
					Proofs: []dto.LoanApprovalValidatorProof{
						{
							ProofURL: "https://example.com/proof1.pdf",
							Category: "identity",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error - create loan approval fails",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("BeginTransaction", context.Background()).Return(&gorm.DB{}, nil)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(&models.Loan{
						ID:     1,
						UUID:   "loan-uuid-123",
						Status: enums.LoanStatusProposed,
					}, nil)
					m.On("CreateLoanApproval", context.Background(), mock.Anything, mock.Anything).Return(errors.New("database error"))
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.ApproveLoanRequest{
					LoanUUID:   "loan-uuid-123",
					EmployeeID: "emp123",
					ApprovedAt: time.Now(),
					Proofs: []dto.LoanApprovalValidatorProof{
						{
							ProofURL: "https://example.com/proof1.pdf",
							Category: "identity",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error - commit fails",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("BeginTransaction", context.Background()).Return(&gorm.DB{}, nil)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(&models.Loan{
						ID:     1,
						UUID:   "loan-uuid-123",
						Status: enums.LoanStatusProposed,
					}, nil)
					m.On("CreateLoanApproval", context.Background(), mock.Anything, mock.Anything).Return(nil)
					m.On("CreateLoanApprovalValidator", context.Background(), mock.Anything, mock.Anything).Return(nil)
					m.On("CreateLoanApprovalValidatorProof", context.Background(), mock.Anything, mock.Anything).Return(nil)
					m.On("UpdateLoan", context.Background(), mock.Anything, mock.Anything, []string{"status"}).Return(nil)
					m.On("Commit", context.Background(), mock.Anything).Return(errors.New("commit failed"))
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.ApproveLoanRequest{
					LoanUUID:   "loan-uuid-123",
					EmployeeID: "emp123",
					ApprovedAt: time.Now(),
					Proofs: []dto.LoanApprovalValidatorProof{
						{
							ProofURL: "https://example.com/proof1.pdf",
							Category: "identity",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LoanService{
				repo:               tt.fields.repo,
				notificationClient: tt.fields.notificationClient,
			}
			if err := s.ApproveLoanWithValidators(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("LoanService.ApproveLoanWithValidators() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoanService_InvestLoan(t *testing.T) {
	type fields struct {
		repo               repository.LoanRepositoryInterface
		notificationClient client.NotificationClientInterface
	}
	type args struct {
		ctx context.Context
		req dto.InvestLoanRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(&models.Loan{
						ID:               1,
						UUID:             "loan-uuid-123",
						Status:           enums.LoanStatusApproved,
						PrincipalAmount:  1000.0,
						InvestmentAmount: 0,
					}, nil)
					m.On("BeginTransaction", context.Background()).Return(&gorm.DB{}, nil)
					m.On("UpdateLoan", context.Background(), mock.Anything, mock.Anything, []string{"investment_amount"}).Return(nil)
					m.On("CreateInvestment", context.Background(), mock.Anything, mock.Anything).Return(nil)
					m.On("Commit", context.Background(), mock.Anything).Return(nil)
					m.On("GetInvestmentsByLoanID", context.Background(), 1).Return([]models.Investment{
						{InvestorID: "investor123", AgreementLetterURL: "https://example.com/agreement.pdf"},
					}, nil)
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					m.On("SendEmail", context.Background(), mock.Anything).Return(nil)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.InvestLoanRequest{
					LoanUUID:   "loan-uuid-123",
					InvestorID: "investor123",
					Amount:     500.0,
				},
			},
			wantErr: false,
		},
		{
			name: "error - loan not found",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(nil, errors.New("loan not found"))
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.InvestLoanRequest{
					LoanUUID:   "loan-uuid-123",
					InvestorID: "investor123",
					Amount:     500.0,
				},
			},
			wantErr: true,
		},
		{
			name: "error - loan not approved",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(&models.Loan{
						ID:     1,
						UUID:   "loan-uuid-123",
						Status: enums.LoanStatusProposed,
					}, nil)
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.InvestLoanRequest{
					LoanUUID:   "loan-uuid-123",
					InvestorID: "investor123",
					Amount:     500.0,
				},
			},
			wantErr: true,
		},
		{
			name: "error - investment amount exceeds principal",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(&models.Loan{
						ID:               1,
						UUID:             "loan-uuid-123",
						Status:           enums.LoanStatusApproved,
						PrincipalAmount:  1000.0,
						InvestmentAmount: 600.0,
					}, nil)
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.InvestLoanRequest{
					LoanUUID:   "loan-uuid-123",
					InvestorID: "investor123",
					Amount:     500.0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LoanService{
				repo:               tt.fields.repo,
				notificationClient: tt.fields.notificationClient,
			}
			if err := s.InvestLoan(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("LoanService.InvestLoan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoanService_CreateLoanDisbursement(t *testing.T) {
	type fields struct {
		repo               repository.LoanRepositoryInterface
		notificationClient client.NotificationClientInterface
	}
	type args struct {
		ctx context.Context
		req dto.CreateLoanDisbursementRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(&models.Loan{
						ID:     1,
						UUID:   "loan-uuid-123",
						Status: enums.LoanStatusInvested,
					}, nil)
					m.On("BeginTransaction", context.Background()).Return(&gorm.DB{}, nil)
					m.On("CreateLoanDisbursement", context.Background(), mock.Anything, mock.Anything).Return(nil)
					m.On("UpdateLoan", context.Background(), mock.Anything, mock.Anything, []string{"status"}).Return(nil)
					m.On("Commit", context.Background(), mock.Anything).Return(nil)
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.CreateLoanDisbursementRequest{
					LoanUUID:                 "loan-uuid-123",
					EmployeeID:               "emp123",
					SignedAgreementLetterURL: "https://example.com/signed-agreement.pdf",
					DisbursedAt:              time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "error - loan not found",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(nil, errors.New("loan not found"))
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.CreateLoanDisbursementRequest{
					LoanUUID:                 "loan-uuid-123",
					EmployeeID:               "emp123",
					SignedAgreementLetterURL: "https://example.com/signed-agreement.pdf",
					DisbursedAt:              time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name: "error - loan not invested",
			fields: fields{
				repo: func() *mocks.LoanRepositoryInterface {
					m := mocks.NewLoanRepositoryInterface(t)
					m.On("GetLoanByUUID", context.Background(), "loan-uuid-123").Return(&models.Loan{
						ID:     1,
						UUID:   "loan-uuid-123",
						Status: enums.LoanStatusApproved,
					}, nil)
					return m
				}(),
				notificationClient: func() *mocks.NotificationClientInterface {
					m := mocks.NewNotificationClientInterface(t)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: dto.CreateLoanDisbursementRequest{
					LoanUUID:                 "loan-uuid-123",
					EmployeeID:               "emp123",
					SignedAgreementLetterURL: "https://example.com/signed-agreement.pdf",
					DisbursedAt:              time.Now(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LoanService{
				repo:               tt.fields.repo,
				notificationClient: tt.fields.notificationClient,
			}
			if err := s.CreateLoanDisbursement(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("LoanService.CreateLoanDisbursement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
