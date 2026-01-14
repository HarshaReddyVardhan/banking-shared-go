package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "SecureP@ssw0rd!"
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Check format
	parts := strings.Split(hash, "$")
	assert.Equal(t, 6, len(parts))
	assert.Equal(t, "argon2id", parts[1])
}

func TestVerifyPassword(t *testing.T) {
	password := "SecureP@ssw0rd!"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	match, err := VerifyPassword(password, hash)
	assert.NoError(t, err)
	assert.True(t, match)

	match, err = VerifyPassword("WrongPassword", hash)
	assert.NoError(t, err)
	assert.False(t, match)
}

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	// Fill key with some data for test
	copy(key, []byte("12345678901234567890123456789012"))

	data := []byte("Sensitive Data")

	encrypted, err := Encrypt(data, key)
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, data, encrypted)

	decrypted, err := Decrypt(encrypted, key)
	assert.NoError(t, err)
	assert.Equal(t, data, decrypted)
}

func TestEncryptInvalidKey(t *testing.T) {
	key := []byte("short")
	data := []byte("data")
	_, err := Encrypt(data, key)
	assert.Error(t, err)
}

func TestDecryptInvalidKey(t *testing.T) {
	key := []byte("short")
	data := []byte("data")
	_, err := Decrypt(data, key)
	assert.Error(t, err)
}

func TestDecryptCorrupted(t *testing.T) {
	key := make([]byte, 32)
	copy(key, []byte("12345678901234567890123456789012"))

	data := []byte("Sensitive Data")
	encrypted, _ := Encrypt(data, key)

	// Corrupt data
	encrypted[len(encrypted)-1] ^= 0xFF

	_, err := Decrypt(encrypted, key)
	assert.Error(t, err)
}
