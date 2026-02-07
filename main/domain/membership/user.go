package membership

import (
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
)

type User struct {
	Id           domain.Id[User]
	FirstName    string
	LastName     string
	BirthDate    time.Time
	Email        EmailAddress
	TaxCode      string // Italian tax code (Codice Fiscale)
	Addresses    []Address
	PhoneNumbers []PhoneNumber
}
