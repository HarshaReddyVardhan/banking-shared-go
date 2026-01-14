package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionStatus_IsFinal(t *testing.T) {
	tests := []struct {
		status TransactionStatus
		want   bool
	}{
		{StatusPending, false},
		{StatusAnalyzing, false},
		{StatusWaitingReview, false},
		{StatusApproved, false},
		{StatusCompleted, true},
		{StatusRejected, true},
		{StatusFailed, true},
		{StatusCancelled, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsFinal())
		})
	}
}

func TestIsValidCurrency(t *testing.T) {
	tests := []struct {
		code string
		want bool
	}{
		{"USD", true},
		{"EUR", true},
		{"GBP", true},
		{"JPY", true},
		{"CAD", true},
		{"AUD", true},
		{"XYZ", false}, // Invalid
		{"", false},
		{"usd", false}, // Case sensitive check in implementation
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidCurrency(tt.code))
		})
	}
}
