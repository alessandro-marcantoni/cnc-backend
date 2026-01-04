package membership

import (
	"regexp"
	"strings"

	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

// EmailAddress is an immutable value object
type EmailAddress struct {
	Value string
}

// NewEmailAddress creates a validated email (smart constructor)
func NewEmailAddress(email string) result.Result[EmailAddress] {
	if email == "" {
		return result.Err[EmailAddress](errors.EmailError{Description: "empty email"})
	}

	email = strings.ToLower(strings.TrimSpace(email))

	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailPattern, email)
	if !matched {
		return result.Err[EmailAddress](errors.EmailError{Description: "invalid email format"})
	}

	return result.Ok(EmailAddress{Value: email})
}

// String implements Stringer interface
func (e EmailAddress) String() string {
	return e.Value
}

// Equals checks equality (value object semantic)
func (e EmailAddress) Equals(other EmailAddress) bool {
	return e.Value == other.Value
}
