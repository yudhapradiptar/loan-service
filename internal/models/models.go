package models

import (
	"time"
)

type Loan struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	UserID       int        `json:"user_id" gorm:"not null"`
	Amount       float64    `json:"amount" gorm:"not null"`
	InterestRate float64    `json:"interest_rate" gorm:"not null"`
	Term         int        `json:"term" gorm:"not null"`
	Status       LoanStatus `json:"status" gorm:"default:1"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateLoanRequest struct {
	UserID       int     `json:"user_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	InterestRate float64 `json:"interest_rate" binding:"required,gt=0"`
	Term         int     `json:"term" binding:"required,gt=0"`
}

type UpdateLoanRequest struct {
	Amount       *float64    `json:"amount,omitempty"`
	InterestRate *float64    `json:"interest_rate,omitempty"`
	Term         *int        `json:"term,omitempty"`
	Status       *LoanStatus `json:"status,omitempty"`
}
