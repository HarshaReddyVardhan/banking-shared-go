package events

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseEvent(t *testing.T) {
	e := NewBaseEvent(EventTypeTransactionInitiated, "test-source")
	assert.NotEmpty(t, e.EventID)
	assert.Equal(t, EventTypeTransactionInitiated, e.EventType)
	assert.False(t, e.Timestamp.IsZero())
	assert.Equal(t, "1.0", e.Version)
	assert.Equal(t, "test-source", e.Source)
}

func TestBaseEvent_WithCorrelation(t *testing.T) {
	e := NewBaseEvent(EventTypeTransactionInitiated, "test")
	e = e.WithCorrelation("corr-123")
	assert.Equal(t, "corr-123", e.CorrelationID)
}

func TestBaseEvent_WithCausation(t *testing.T) {
	e := NewBaseEvent(EventTypeTransactionInitiated, "test")
	e = e.WithCausation("cause-123")
	assert.Equal(t, "cause-123", e.CausationID)
}

func TestTransactionInitiatedEvent_Marshal(t *testing.T) {
	uid := uuid.New()
	txID := uuid.New()

	e := TransactionInitiatedEvent{
		BaseEvent:     NewBaseEvent(EventTypeTransactionInitiated, "test"),
		TransactionID: txID,
		UserID:        uid,
		Amount:        decimal.NewFromFloat(100.50),
		Currency:      "USD",
		TransferType:  "P2P",
	}

	data, err := json.Marshal(&e)
	assert.NoError(t, err)

	var decoded TransactionInitiatedEvent
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, uid, decoded.UserID)
	assert.Equal(t, txID, decoded.TransactionID)
	assert.True(t, e.Amount.Equal(decoded.Amount))
	assert.Equal(t, "100.5", decoded.Amount.String())
}

func TestEventKeys(t *testing.T) {
	uid := uuid.New()
	txID := uuid.New()

	t.Run("TransactionInitiatedEvent", func(t *testing.T) {
		e := &TransactionInitiatedEvent{UserID: uid}
		assert.Equal(t, uid.String(), e.Key())
	})

	t.Run("FraudAnalysisCompleteEvent", func(t *testing.T) {
		e := &FraudAnalysisCompleteEvent{TransactionID: txID}
		assert.Equal(t, txID.String(), e.Key())
	})

	t.Run("TransactionCompletedEvent", func(t *testing.T) {
		e := &TransactionCompletedEvent{TransactionID: txID}
		assert.Equal(t, txID.String(), e.Key())
	})

	t.Run("UserCreatedEvent", func(t *testing.T) {
		e := &UserCreatedEvent{UserID: uid}
		assert.Equal(t, uid.String(), e.Key())
	})

	t.Run("AuditLogEvent", func(t *testing.T) {
		e := &AuditLogEvent{ActorID: "user-123"}
		assert.Equal(t, "user-123", e.Key())
	})
}
