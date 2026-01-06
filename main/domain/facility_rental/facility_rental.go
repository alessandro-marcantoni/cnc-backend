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

func (s SimpleRentedFacility) GetId() domain.Id[RentedFacility] {
	return s.Id
}

func (s SimpleRentedFacility) GetMemberId() domain.Id[membership.Member] {
	return s.MemberId
}

func (s SimpleRentedFacility) GetFacility() domain.Id[Facility] {
	return s.Facility
}

func (s SimpleRentedFacility) GetValidity() RentalValidity {
	return s.Validity
}

func (s SimpleRentedFacility) GetPayment() payment.Payment {
	return s.Payment
}

func (s SimpleRentedFacility) GetType() RentedFacilityType {
	return SimpleFacility
}

func (r RentedFacilityWithBoat) GetId() domain.Id[RentedFacility] {
	return r.Id
}

func (r RentedFacilityWithBoat) GetMemberId() domain.Id[membership.Member] {
	return r.MemberId
}

func (r RentedFacilityWithBoat) GetFacility() domain.Id[Facility] {
	return r.Facility
}

func (r RentedFacilityWithBoat) GetValidity() RentalValidity {
	return r.Validity
}

func (r RentedFacilityWithBoat) GetPayment() payment.Payment {
	return r.Payment
}

func (r RentedFacilityWithBoat) GetType() RentedFacilityType {
	return BoatFacility
}
