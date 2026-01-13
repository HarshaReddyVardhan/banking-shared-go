// Package models provides shared domain models for banking services.
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TransactionStatus represents the current state of a transaction
type TransactionStatus string

const (
	StatusPending       TransactionStatus = "PENDING"
	StatusAnalyzing     TransactionStatus = "ANALYZING"
	StatusWaitingReview TransactionStatus = "WAITING_REVIEW"
	StatusApproved      TransactionStatus = "APPROVED"
	StatusRejected      TransactionStatus = "REJECTED"
	StatusCompleted     TransactionStatus = "COMPLETED"
	StatusFailed        TransactionStatus = "FAILED"
	StatusCancelled     TransactionStatus = "CANCELLED"
)

// IsFinal returns true if the status is a terminal state
func (s TransactionStatus) IsFinal() bool {
	return s == StatusCompleted || s == StatusRejected || s == StatusFailed || s == StatusCancelled
}

// FraudDecision represents the decision from fraud analysis
type FraudDecision string

const (
	FraudDecisionApproved   FraudDecision = "APPROVED"
	FraudDecisionSuspicious FraudDecision = "SUSPICIOUS"
	FraudDecisionRejected   FraudDecision = "REJECTED"
)

// TransferType represents the type of transfer
type TransferType string

const (
	TransferTypeInternal TransferType = "INTERNAL"
	TransferTypeExternal TransferType = "EXTERNAL"
	TransferTypeWire     TransferType = "WIRE"
	TransferTypeACH      TransferType = "ACH"
)

// InitiationMethod represents how the transfer was initiated
type InitiationMethod string

const (
	InitiationWeb    InitiationMethod = "WEB"
	InitiationMobile InitiationMethod = "MOBILE"
	InitiationAPI    InitiationMethod = "API"
	InitiationBranch InitiationMethod = "BRANCH"
)

// Currency represents supported currencies
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
	CurrencyJPY Currency = "JPY"
	CurrencyCAD Currency = "CAD"
	CurrencyAUD Currency = "AUD"
)

// SupportedCurrencies returns a list of supported currencies
func SupportedCurrencies() []Currency {
	return []Currency{CurrencyUSD, CurrencyEUR, CurrencyGBP, CurrencyJPY, CurrencyCAD, CurrencyAUD}
}

// IsValidCurrency checks if a currency code is supported
func IsValidCurrency(code string) bool {
	for _, c := range SupportedCurrencies() {
		if string(c) == code {
			return true
		}
	}
	return false
}

// Transaction represents the core transaction model shared across services
type Transaction struct {
	ID               uuid.UUID         `json:"id" db:"id"`
	UserID           uuid.UUID         `json:"user_id" db:"user_id"`
	FromAccountID    uuid.UUID         `json:"from_account_id" db:"from_account_id"`
	ToAccountID      uuid.UUID         `json:"to_account_id" db:"to_account_id"`
	Amount           decimal.Decimal   `json:"amount" db:"amount"`
	Currency         Currency          `json:"currency" db:"currency"`
	Status           TransactionStatus `json:"status" db:"status"`
	TransferType     TransferType      `json:"transfer_type" db:"transfer_type"`
	InitiationMethod InitiationMethod  `json:"initiation_method" db:"initiation_method"`
	Reference        string            `json:"reference" db:"reference"`
	Memo             string            `json:"memo,omitempty" db:"memo"`
	FraudScore       *float64          `json:"fraud_score,omitempty" db:"fraud_score"`
	FraudDecision    *FraudDecision    `json:"fraud_decision,omitempty" db:"fraud_decision"`
	SourceIP         string            `json:"source_ip,omitempty" db:"source_ip"`
	DeviceID         string            `json:"device_id,omitempty" db:"device_id"`
	SessionID        string            `json:"session_id,omitempty" db:"session_id"`
	UserAgent        string            `json:"user_agent,omitempty" db:"user_agent"`
	InitiatedAt      time.Time         `json:"initiated_at" db:"initiated_at"`
	CompletedAt      *time.Time        `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at" db:"updated_at"`
}

// User represents the core user model shared across services
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	PhoneNumber  string    `json:"phone_number,omitempty" db:"phone_number"`
	Tier         string    `json:"tier" db:"tier"` // BASIC, PREMIUM, ENTERPRISE
	Status       string    `json:"status" db:"status"`
	RiskScore    float64   `json:"risk_score" db:"risk_score"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	IsMFAEnabled bool      `json:"is_mfa_enabled" db:"is_mfa_enabled"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Account represents a bank account shared across services
type Account struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	UserID           uuid.UUID       `json:"user_id" db:"user_id"`
	AccountNumber    string          `json:"account_number" db:"account_number"`
	AccountType      string          `json:"account_type" db:"account_type"` // CHECKING, SAVINGS
	Currency         Currency        `json:"currency" db:"currency"`
	Balance          decimal.Decimal `json:"balance" db:"balance"`
	AvailableBalance decimal.Decimal `json:"available_balance" db:"available_balance"`
	Status           string          `json:"status" db:"status"` // ACTIVE, FROZEN, CLOSED
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}
