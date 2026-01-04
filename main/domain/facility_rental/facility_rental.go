package facilityrental

import (
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
)

type RentedFacility interface {
	GetId() domain.Id[RentedFacility]
	GetMemberId() domain.Id[membership.Member]
	GetFacility() domain.Id[Facility]
	GetValidity() RentalValidity
	GetPayment() payment.Payment
	GetType() RentedFacilityType
}

type RentedFacilityType string

const (
	SimpleFacility RentedFacilityType = "SIMPLE_FACILITY"
	BoatFacility   RentedFacilityType = "BOAT_FACILITY"
)

type SimpleRentedFacility struct {
	Id       domain.Id[RentedFacility]
	MemberId domain.Id[membership.Member]
	Facility domain.Id[Facility]
	Validity RentalValidity
	Payment  payment.Payment
}

type RentedFacilityWithBoat struct {
	Id       domain.Id[RentedFacility]
	MemberId domain.Id[membership.Member]
	Facility domain.Id[Facility]
	Validity RentalValidity
	Payment  payment.Payment
	BoatInfo BoatInfo
}

type RentalValidity struct {
	ToDate time.Time
}
