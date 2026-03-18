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
	GetFacility() Facility
	GetValidity() RentalValidity
	GetPayment() payment.Payment
	GetType() RentedFacilityType
	GetPrice() float64
	GetDiscountApplied() bool
}

type RentedFacilityType string

const (
	SimpleFacility    RentedFacilityType = "SIMPLE_FACILITY"
	BoatFacility      RentedFacilityType = "BOAT_FACILITY"
	LeerboardFacility RentedFacilityType = "LEERBOARD_FACILITY"
)

type SimpleRentedFacility struct {
	Id              domain.Id[RentedFacility]
	MemberId        domain.Id[membership.Member]
	Facility        Facility
	Validity        RentalValidity
	Price           float64
	Payment         payment.Payment
	DiscountApplied bool
}

type RentedFacilityWithBoat struct {
	Id              domain.Id[RentedFacility]
	MemberId        domain.Id[membership.Member]
	Facility        Facility
	Validity        RentalValidity
	Price           float64
	Payment         payment.Payment
	BoatInfo        BoatInfo
	DiscountApplied bool
}

type RentedFacilityWithLeerboard struct {
	Id              domain.Id[RentedFacility]
	MemberId        domain.Id[membership.Member]
	Facility        Facility
	Validity        RentalValidity
	Price           float64
	Payment         payment.Payment
	LeerboardInfo   LeerboardInfo
	DiscountApplied bool
}

type RentalValidity struct {
	FromDate time.Time
	ToDate   time.Time
}

func (s SimpleRentedFacility) GetId() domain.Id[RentedFacility] {
	return s.Id
}

func (s SimpleRentedFacility) GetMemberId() domain.Id[membership.Member] {
	return s.MemberId
}

func (s SimpleRentedFacility) GetFacility() Facility {
	return s.Facility
}

func (s SimpleRentedFacility) GetValidity() RentalValidity {
	return s.Validity
}

func (s SimpleRentedFacility) GetPrice() float64 {
	return s.Price
}

func (s SimpleRentedFacility) GetPayment() payment.Payment {
	return s.Payment
}

func (s SimpleRentedFacility) GetType() RentedFacilityType {
	return SimpleFacility
}

func (s SimpleRentedFacility) GetDiscountApplied() bool {
	return s.DiscountApplied
}

func (r RentedFacilityWithBoat) GetId() domain.Id[RentedFacility] {
	return r.Id
}

func (r RentedFacilityWithBoat) GetMemberId() domain.Id[membership.Member] {
	return r.MemberId
}

func (r RentedFacilityWithBoat) GetFacility() Facility {
	return r.Facility
}

func (r RentedFacilityWithBoat) GetValidity() RentalValidity {
	return r.Validity
}

func (r RentedFacilityWithBoat) GetPrice() float64 {
	return r.Price
}

func (r RentedFacilityWithBoat) GetPayment() payment.Payment {
	return r.Payment
}

func (r RentedFacilityWithBoat) GetType() RentedFacilityType {
	return BoatFacility
}

func (r RentedFacilityWithBoat) GetDiscountApplied() bool {
	return r.DiscountApplied
}

func (r RentedFacilityWithLeerboard) GetId() domain.Id[RentedFacility] {
	return r.Id
}

func (r RentedFacilityWithLeerboard) GetMemberId() domain.Id[membership.Member] {
	return r.MemberId
}

func (r RentedFacilityWithLeerboard) GetFacility() Facility {
	return r.Facility
}

func (r RentedFacilityWithLeerboard) GetValidity() RentalValidity {
	return r.Validity
}

func (r RentedFacilityWithLeerboard) GetPrice() float64 {
	return r.Price
}

func (r RentedFacilityWithLeerboard) GetPayment() payment.Payment {
	return r.Payment
}

func (r RentedFacilityWithLeerboard) GetType() RentedFacilityType {
	return LeerboardFacility
}

func (r RentedFacilityWithLeerboard) GetDiscountApplied() bool {
	return r.DiscountApplied
}
