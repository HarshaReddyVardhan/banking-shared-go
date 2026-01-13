// Package events provides standardized Kafka event definitions for banking services.
package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// EventType represents the type of banking event
type EventType string

// Standard event types used across all banking services
const (
	// Transaction Events
	EventTypeTransactionInitiated     EventType = "TransactionInitiated"
	EventTypeTransactionAnalyzing     EventType = "TransactionAnalyzing"
	EventTypeTransactionApproved      EventType = "TransactionApproved"
	EventTypeTransactionRejected      EventType = "TransactionRejected"
	EventTypeTransactionCompleted     EventType = "TransactionCompleted"
	EventTypeTransactionFailed        EventType = "TransactionFailed"
	EventTypeTransactionCancelled     EventType = "TransactionCancelled"
	EventTypeTransactionWaitingReview EventType = "TransactionWaitingReview"

	// Fraud Events
	EventTypeFraudAnalysisComplete EventType = "FraudAnalysisComplete"
	EventTypeFraudSuspected        EventType = "FraudSuspected"
	EventTypeFraudReviewComplete   EventType = "FraudReviewComplete"
	EventTypeManualReviewRequired  EventType = "ManualReviewRequired"
	EventTypeBlocklistMatch        EventType = "BlocklistMatch"

	// User Events
	EventTypeUserCreated         EventType = "UserCreated"
	EventTypeUserUpdated         EventType = "UserUpdated"
	EventTypeUserLocked          EventType = "UserLocked"
	EventTypeUserPasswordChanged EventType = "UserPasswordChanged"

	// Auth Events
	EventTypeLoginSuccess  EventType = "LoginSuccess"
	EventTypeLoginFailed   EventType = "LoginFailed"
	EventTypeMFAEnabled    EventType = "MFAEnabled"
	EventTypeTokenRevoked  EventType = "TokenRevoked"
	EventTypeJWTKeyRotated EventType = "JWTKeyRotated"
	EventTypeSecurityAlert EventType = "SecurityAlert"

	// Notification Events
	EventTypeNotificationSent   EventType = "NotificationSent"
	EventTypeNotificationFailed EventType = "NotificationFailed"

	// AML Events
	EventTypeAMLScreeningComplete EventType = "AMLScreeningComplete"
	EventTypeSARFiled             EventType = "SARFiled"
	EventTypeRiskProfileUpdated   EventType = "RiskProfileUpdated"

	// Audit Events
	EventTypeAuditLogCreated EventType = "AuditLogCreated"
)

// BaseEvent contains common fields for all events
type BaseEvent struct {
	EventID       uuid.UUID `json:"event_id"`
	EventType     EventType `json:"event_type"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	CausationID   string    `json:"causation_id,omitempty"`
	Source        string    `json:"source"`
}

// NewBaseEvent creates a new base event with default values
func NewBaseEvent(eventType EventType, source string) BaseEvent {
	return BaseEvent{
		EventID:   uuid.New(),
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
		Source:    source,
	}
}

// WithCorrelation sets the correlation ID and returns the event
func (e BaseEvent) WithCorrelation(correlationID string) BaseEvent {
	e.CorrelationID = correlationID
	return e
}

// WithCausation sets the causation ID and returns the event
func (e BaseEvent) WithCausation(causationID string) BaseEvent {
	e.CausationID = causationID
	return e
}

// EventMetadata contains context about the event source
type EventMetadata struct {
	SourceIP         string `json:"source_ip,omitempty"`
	UserAgent        string `json:"user_agent,omitempty"`
	DeviceID         string `json:"device_id,omitempty"`
	SessionID        string `json:"session_id,omitempty"`
	InitiationMethod string `json:"initiation_method,omitempty"`
}

// TransactionInitiatedEvent is published when a transfer is initiated
type TransactionInitiatedEvent struct {
	BaseEvent
	TransactionID uuid.UUID       `json:"transaction_id"`
	UserID        uuid.UUID       `json:"user_id"`
	FromAccountID uuid.UUID       `json:"from_account_id"`
	ToAccountID   uuid.UUID       `json:"to_account_id"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	TransferType  string          `json:"transfer_type"`
	Memo          string          `json:"memo,omitempty"`
	Metadata      EventMetadata   `json:"metadata"`
}

// Key returns the partition key for Kafka (user_id for consistent ordering)
func (e *TransactionInitiatedEvent) Key() string {
	return e.UserID.String()
}

// MarshalJSON serializes the event to JSON
func (e *TransactionInitiatedEvent) MarshalJSON() ([]byte, error) {
	type Alias TransactionInitiatedEvent
	return json.Marshal((*Alias)(e))
}

// FraudAnalysisCompleteEvent is published when fraud analysis is complete
type FraudAnalysisCompleteEvent struct {
	BaseEvent
	TransactionID uuid.UUID `json:"transaction_id"`
	UserID        uuid.UUID `json:"user_id"`
	AnalysisID    string    `json:"analysis_id"`
	RiskScore     float64   `json:"risk_score"`
	Decision      string    `json:"decision"` // APPROVED, REJECTED, REVIEW_REQUIRED
	Reasons       []string  `json:"reasons,omitempty"`
	ProcessingMs  int64     `json:"processing_ms"`
}

// Key returns the partition key for Kafka
func (e *FraudAnalysisCompleteEvent) Key() string {
	return e.TransactionID.String()
}

// TransactionCompletedEvent is published when a transaction is completed
type TransactionCompletedEvent struct {
	BaseEvent
	TransactionID  uuid.UUID       `json:"transaction_id"`
	UserID         uuid.UUID       `json:"user_id"`
	Amount         decimal.Decimal `json:"amount"`
	Currency       string          `json:"currency"`
	ProcessingTime int64           `json:"processing_time_ms"`
}

// Key returns the partition key for Kafka
func (e *TransactionCompletedEvent) Key() string {
	return e.TransactionID.String()
}

// UserCreatedEvent is published when a new user is registered
type UserCreatedEvent struct {
	BaseEvent
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Tier      string    `json:"tier"`
}

// Key returns the partition key for Kafka
func (e *UserCreatedEvent) Key() string {
	return e.UserID.String()
}

// AuditLogEvent is published for audit trail
type AuditLogEvent struct {
	BaseEvent
	ActorID      string                 `json:"actor_id"`
	ActorType    string                 `json:"actor_type"` // user, system, admin
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   string                 `json:"resource_id"`
	Details      map[string]interface{} `json:"details,omitempty"`
	IPAddress    string                 `json:"ip_address,omitempty"`
}

// Key returns the partition key for Kafka
func (e *AuditLogEvent) Key() string {
	return e.ActorID
}
