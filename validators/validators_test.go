package validators

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestValidateTransferAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  decimal.Decimal
		wantErr bool
	}{
		{"Valid amount", decimal.NewFromFloat(100.00), false},
		{"Zero amount", decimal.Zero, true},
		{"Negative amount", decimal.NewFromFloat(-10.00), true},
		{"Exceeds limit", MaxTransferAmount.Add(decimal.NewFromFloat(1)), true},
		{"Too many decimals", decimal.NewFromFloat(100.001), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransferAmount(tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAccountNumber(t *testing.T) {
	tests := []struct {
		name    string
		accNum  string
		wantErr bool
	}{
		{"Valid 10 digits", "1234567890", false},
		{"Valid 8 digits", "12345678", false},
		{"Valid 12 digits", "123456789012", false},
		{"Too short", "1234567", true},
		{"Too long", "1234567890123", true},
		{"Non-digits", "12345abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAccountNumber(tt.accNum)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRoutingNumber(t *testing.T) {
	tests := []struct {
		name    string
		routing string
		wantErr bool
	}{
		{"Valid routing", "111000025", false},      // Chase Bank (example, checksum valid)
		{"Valid routing test", "091000022", false}, // Wells Fargo
		{"Invalid length", "12345678", true},
		{"Non-digits", "12345678a", true},
		{"Invalid checksum", "111000026", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRoutingNumber(tt.routing)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid email", "test@example.com", false},
		{"Valid email with dot", "test.user@example.co.uk", false},
		{"Missing @", "testexample.com", true},
		{"Missing domain", "test@", true},
		{"Empty", "", true},
		{"Whitespace", "  test@example.com  ", false}, // Should be trimmed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", "StrongP@ssw0rd!", false},
		{"Too short", "Weak1!", true},
		{"No lower", "STRONGP@SSW0RD!", true},
		{"No upper", "strongp@ssw0rd!", true},
		{"No digit", "StrongP@ssword!", true},
		{"No special", "StrongPassw0rd", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal", "hello world", "hello world"},
		{"Trim", "  hello world  ", "hello world"},
		{"Remove null", "hello\x00world", "helloworld"},
		{"Allow newline", "hello\nworld", "hello\nworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeInput(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
