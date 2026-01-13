// Package validators provides input validation for banking operations.
package validators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

// ValidationError represents a validation failure
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// MaxTransferAmount is the default maximum transfer limit
var MaxTransferAmount = decimal.NewFromInt(1000000)

// ValidateTransferAmount validates a transfer amount
func ValidateTransferAmount(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return ValidationError{Field: "amount", Message: "Amount must be greater than zero"}
	}
	if amount.GreaterThan(MaxTransferAmount) {
		return ValidationError{Field: "amount", Message: fmt.Sprintf("Amount exceeds maximum limit of %s", MaxTransferAmount)}
	}
	// Check decimal places (max 2 for most currencies)
	if amount.Exponent() < -2 {
		return ValidationError{Field: "amount", Message: "Amount has too many decimal places (max 2)"}
	}
	return nil
}

// ValidateAccountNumber validates a bank account number
func ValidateAccountNumber(accountNumber string) error {
	if len(accountNumber) < 8 || len(accountNumber) > 12 {
		return ValidationError{Field: "account_number", Message: "Account number must be between 8 and 12 digits"}
	}
	if !regexp.MustCompile(`^\d+$`).MatchString(accountNumber) {
		return ValidationError{Field: "account_number", Message: "Account number must contain only digits"}
	}
	return nil
}

// ValidateRoutingNumber validates a US ABA routing number (9 digits with checksum)
func ValidateRoutingNumber(routingNumber string) error {
	if !regexp.MustCompile(`^\d{9}$`).MatchString(routingNumber) {
		return ValidationError{Field: "routing_number", Message: "Routing number must be exactly 9 digits"}
	}

	// ABA checksum algorithm
	digits := make([]int, 9)
	for i, c := range routingNumber {
		digits[i] = int(c - '0')
	}

	checksum := (3*(digits[0]+digits[3]+digits[6]) +
		7*(digits[1]+digits[4]+digits[7]) +
		(digits[2] + digits[5] + digits[8])) % 10

	if checksum != 0 {
		return ValidationError{Field: "routing_number", Message: "Invalid routing number checksum"}
	}
	return nil
}

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return ValidationError{Field: "email", Message: "Email is required"}
	}
	// Basic email regex - for production use a more comprehensive validator
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ValidationError{Field: "email", Message: "Invalid email format"}
	}
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return ValidationError{Field: "password", Message: "Password must be at least 12 characters"}
	}
	if len(password) > 128 {
		return ValidationError{Field: "password", Message: "Password must be at most 128 characters"}
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return ValidationError{
			Field:   "password",
			Message: "Password must contain uppercase, lowercase, digit, and special character",
		}
	}
	return nil
}

// SanitizeInput removes potentially dangerous characters from user input
func SanitizeInput(input string) string {
	// Remove null bytes and other control characters
	sanitized := strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, input)

	// Trim whitespace
	return strings.TrimSpace(sanitized)
}
