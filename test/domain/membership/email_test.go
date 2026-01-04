package membership_test

import (
	"testing"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/stretchr/testify/assert"
)

func TestNewEmailAddress_ValidEmail(t *testing.T) {
	// Arrange
	test_email := "test@example.com"

	// Act
	result := membership.NewEmailAddress(test_email)

	// Assert
	assert.True(t, result.IsSuccess())
	assert.Equal(t, "test@example.com", result.Value().Value)
}

func TestNewEmailAddress_InvalidFormat(t *testing.T) {
	// Arrange
	testCases := []struct {
		name  string
		email string
	}{
		{"empty", ""},
		{"no at sign", "invalid-email"},
		{"no domain", "test@"},
		{"no local part", "@example.com"},
		{"spaces", "test @example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := membership.NewEmailAddress(tc.email)

			// Assert
			assert.False(t, result.IsSuccess())
			assert.NotNil(t, result.Error())
		})
	}
}

func TestEmailAddress_NormalizesToLowercase(t *testing.T) {
	// Arrange
	test_email := "Test@EXAMPLE.COM"

	// Act
	result := membership.NewEmailAddress(test_email)

	// Assert
	assert.True(t, result.IsSuccess())
	assert.Equal(t, "test@example.com", result.Value().Value)
}

func TestEmailAddress_Equals(t *testing.T) {
	// Arrange
	email1 := membership.NewEmailAddress("test@example.com").Value()
	email2 := membership.NewEmailAddress("test@example.com").Value()
	email3 := membership.NewEmailAddress("other@example.com").Value()

	// Assert
	assert.True(t, email1.Equals(email2))
	assert.False(t, email1.Equals(email3))
}
