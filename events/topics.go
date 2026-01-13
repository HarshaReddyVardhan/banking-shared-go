// Package events provides topic configuration for Kafka messaging.
package events

// TopicConfig holds Kafka topic names for the banking platform
type TopicConfig struct {
	// Transaction topics
	TransactionInitiated string
	TransactionApproved  string
	TransactionRejected  string
	TransactionCompleted string

	// Fraud topics
	FraudAnalysis  string
	FraudSuspected string
	ManualReview   string

	// User topics
	UserEvents string

	// Auth/Security topics
	SecurityEvents string

	// Notification topics
	Notifications string

	// AML topics
	AMLScreening string
	SARFiling    string

	// Audit topics
	AuditLog string
}

// DefaultTopicConfig returns the default topic configuration
func DefaultTopicConfig() TopicConfig {
	return TopicConfig{
		TransactionInitiated: "banking.transactions.initiated",
		TransactionApproved:  "banking.transactions.approved",
		TransactionRejected:  "banking.transactions.rejected",
		TransactionCompleted: "banking.transactions.completed",

		FraudAnalysis:  "banking.fraud.analysis",
		FraudSuspected: "banking.fraud.suspected",
		ManualReview:   "banking.fraud.manual-review",

		UserEvents: "banking.users.events",

		SecurityEvents: "banking.security.events",

		Notifications: "banking.notifications",

		AMLScreening: "banking.aml.screening",
		SARFiling:    "banking.aml.sar-filing",

		AuditLog: "banking.audit.log",
	}
}

// Topic returns the appropriate Kafka topic for an event type
func (e EventType) Topic(cfg TopicConfig) string {
	switch e {
	case EventTypeTransactionInitiated:
		return cfg.TransactionInitiated
	case EventTypeTransactionApproved:
		return cfg.TransactionApproved
	case EventTypeTransactionRejected:
		return cfg.TransactionRejected
	case EventTypeTransactionCompleted:
		return cfg.TransactionCompleted
	case EventTypeFraudAnalysisComplete, EventTypeFraudReviewComplete:
		return cfg.FraudAnalysis
	case EventTypeFraudSuspected, EventTypeBlocklistMatch:
		return cfg.FraudSuspected
	case EventTypeManualReviewRequired:
		return cfg.ManualReview
	case EventTypeUserCreated, EventTypeUserUpdated, EventTypeUserLocked:
		return cfg.UserEvents
	case EventTypeLoginSuccess, EventTypeLoginFailed, EventTypeSecurityAlert, EventTypeJWTKeyRotated:
		return cfg.SecurityEvents
	case EventTypeNotificationSent, EventTypeNotificationFailed:
		return cfg.Notifications
	case EventTypeAMLScreeningComplete, EventTypeRiskProfileUpdated:
		return cfg.AMLScreening
	case EventTypeSARFiled:
		return cfg.SARFiling
	case EventTypeAuditLogCreated:
		return cfg.AuditLog
	default:
		return cfg.AuditLog // Default to audit log for unknown events
	}
}
