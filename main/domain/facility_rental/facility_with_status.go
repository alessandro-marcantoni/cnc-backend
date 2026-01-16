package facilityrental

import (
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
)

type FacilityWithStatus struct {
	Id                      domain.Id[Facility]
	FacilityTypeId          domain.Id[FacilityType]
	Identifier              string
	FacilityTypeName        FacilityName
	FacilityTypeDescription string
	SuggestedPrice          float64
	IsRented                bool
	ExpiresAt               *time.Time
	RentedByMemberId        *int64
	RentedByMemberFirstName *string
	RentedByMemberLastName  *string
}
