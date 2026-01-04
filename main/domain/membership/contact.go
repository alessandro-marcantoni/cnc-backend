package membership

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type PhoneNumber struct {
	Prefix *string
	Number string
}

// NewPhoneNumber creates a validated phone number (smart constructor)
func NewPhoneNumber(prefix string, number string) result.Result[PhoneNumber] {
	if len(number) < 10 || len(number) > 15 {
		return result.Err[PhoneNumber](errors.PhoneNumberError{Description: "invalid phone number format"})
	}
	if prefix == "" {
		return result.Ok(PhoneNumber{nil, number})
	}
	if len(prefix) != 3 {
		return result.Err[PhoneNumber](errors.PhoneNumberError{Description: "invalid prefix format"})
	}
	return result.Ok(PhoneNumber{&prefix, number})
}

func (p PhoneNumber) String() string {
	if p.Prefix == nil {
		return p.Number
	}
	return *p.Prefix + p.Number
}
