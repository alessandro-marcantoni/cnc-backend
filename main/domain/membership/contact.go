package membership

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type PhoneNumber struct {
	Number string
}

// NewPhoneNumber creates a validated phone number (smart constructor)
func NewPhoneNumber(number string) result.Result[PhoneNumber] {
	if len(number) < 10 || len(number) > 15 {
		return result.Err[PhoneNumber](errors.PhoneNumberError{Description: "invalid phone number format"})
	}
	return result.Ok(PhoneNumber{Number: number})
}

func (p PhoneNumber) String() string {
	return p.Number
}
