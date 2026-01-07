package persistence

import (
	"database/sql"
	_ "embed"
	"errors"
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	facilityrental "github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

//go:embed queries/get_rented_facilities_by_member.sql
var getRentedFacilitiesByMemberQuery string

//go:embed queries/get_facilities_catalog.sql
var getFacilitiesCatalogQuery string

//go:embed queries/get_facilities_by_type.sql
var getFacilitiesByTypeQuery string

type SQLFacilityRepository struct {
	db *sql.DB
}

func NewSQLFacilityRepository(db *sql.DB) *SQLFacilityRepository {
	return &SQLFacilityRepository{db: db}
}

func (r *SQLFacilityRepository) GetFacilitiesCatalog() []facilityrental.FacilityType {
	rows, err := r.db.Query(getFacilitiesCatalogQuery)
	if err != nil {
		return []facilityrental.FacilityType{}
	}
	defer rows.Close()

	var facilityTypes []facilityrental.FacilityType

	for rows.Next() {
		var id int64
		var name string
		var description sql.NullString
		var suggestedPrice float64

		err := rows.Scan(&id, &name, &description, &suggestedPrice)
		if err != nil {
			continue
		}

		facilityType := facilityrental.FacilityType{
			Id:             domain.Id[facilityrental.FacilityType]{Value: id},
			FacilityName:   facilityrental.ToFacilityName(name),
			Description:    description.String,
			SuggestedPrice: suggestedPrice,
		}
		facilityTypes = append(facilityTypes, facilityType)
	}

	return facilityTypes
}

func (r *SQLFacilityRepository) GetFacilitiesByType(facilityTypeId domain.Id[facilityrental.FacilityType]) []facilityrental.FacilityWithStatus {
	rows, err := r.db.Query(getFacilitiesByTypeQuery, facilityTypeId.Value)
	if err != nil {
		return []facilityrental.FacilityWithStatus{}
	}
	defer rows.Close()

	var facilities []facilityrental.FacilityWithStatus

	for rows.Next() {
		var id int64
		var facilityTypeId int64
		var identifier string
		var facilityTypeName string
		var facilityTypeDescription sql.NullString
		var suggestedPrice float64
		var isRented bool
		var expiresAt sql.NullTime

		err := rows.Scan(
			&id,
			&facilityTypeId,
			&identifier,
			&facilityTypeName,
			&facilityTypeDescription,
			&suggestedPrice,
			&isRented,
			&expiresAt,
		)
		if err != nil {
			continue
		}

		var expiresAtPtr *time.Time
		if expiresAt.Valid {
			expiresAtPtr = &expiresAt.Time
		}

		facility := facilityrental.FacilityWithStatus{
			Id:                      domain.Id[facilityrental.Facility]{Value: id},
			FacilityTypeId:          domain.Id[facilityrental.FacilityType]{Value: facilityTypeId},
			Identifier:              identifier,
			FacilityTypeName:        facilityrental.ToFacilityName(facilityTypeName),
			FacilityTypeDescription: facilityTypeDescription.String,
			SuggestedPrice:          suggestedPrice,
			IsRented:                isRented,
			ExpiresAt:               expiresAtPtr,
		}
		facilities = append(facilities, facility)
	}

	return facilities
}

func (r *SQLFacilityRepository) GetAvailableFacilities(serviceType facilityrental.FacilityName) []facilityrental.Facility {
	// TODO: Implement this method
	return []facilityrental.Facility{}
}

func (r *SQLFacilityRepository) GetFacilitiesRentedByMember(memberId domain.Id[membership.User]) []facilityrental.RentedFacility {
	rows, err := r.db.Query(getRentedFacilitiesByMemberQuery, memberId.Value)
	if err != nil {
		return []facilityrental.RentedFacility{}
	}
	defer rows.Close()

	var rentedFacilities []facilityrental.RentedFacility

	for rows.Next() {
		var dto GetRentedFacilitiesByMemberQueryResult
		err := rows.Scan(
			&dto.RentedFacilityID,
			&dto.RentedAt,
			&dto.ExpiresAt,
			&dto.FacilityID,
			&dto.FacilityIdentifier,
			&dto.FacilityTypeID,
			&dto.FacilityType,
			&dto.FacilityTypeDesc,
			&dto.SuggestedPrice,
			&dto.BoatID,
			&dto.BoatName,
			&dto.LengthMeters,
			&dto.WidthMeters,
			&dto.PaymentID,
			&dto.PaymentAmount,
			&dto.PaymentCurrency,
			&dto.PaymentPaidAt,
			&dto.PaymentMethod,
			&dto.TransactionRef,
		)
		if err != nil {
			continue
		}

		rentedFacility := convertDTOToRentedFacility(dto)
		rentedFacilities = append(rentedFacilities, rentedFacility)
	}

	return rentedFacilities
}

func (r *SQLFacilityRepository) RentFacility(memberId domain.Id[membership.User], facilityId domain.Id[facilityrental.Facility], boatInfo *domain.Id[facilityrental.BoatInfo]) result.Result[facilityrental.RentedFacility] {
	// TODO: Implement this method
	return result.Err[facilityrental.RentedFacility](errors.New("method not implemented"))
}

func (r *SQLFacilityRepository) GetRentedFacilityDTOs(memberId int64) ([]GetRentedFacilitiesByMemberQueryResult, error) {
	rows, err := r.db.Query(getRentedFacilitiesByMemberQuery, memberId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dtos []GetRentedFacilitiesByMemberQueryResult

	for rows.Next() {
		var dto GetRentedFacilitiesByMemberQueryResult
		err := rows.Scan(
			&dto.RentedFacilityID,
			&dto.RentedAt,
			&dto.ExpiresAt,
			&dto.FacilityID,
			&dto.FacilityIdentifier,
			&dto.FacilityTypeID,
			&dto.FacilityType,
			&dto.FacilityTypeDesc,
			&dto.SuggestedPrice,
			&dto.BoatID,
			&dto.BoatName,
			&dto.LengthMeters,
			&dto.WidthMeters,
			&dto.PaymentID,
			&dto.PaymentAmount,
			&dto.PaymentCurrency,
			&dto.PaymentPaidAt,
			&dto.PaymentMethod,
			&dto.TransactionRef,
		)
		if err != nil {
			continue
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

func convertDTOToRentedFacility(dto GetRentedFacilitiesByMemberQueryResult) facilityrental.RentedFacility {
	facilityType := facilityrental.FacilityType{
		FacilityName:   facilityrental.ToFacilityName(dto.FacilityType),
		SuggestedPrice: dto.SuggestedPrice,
	}

	facility := facilityrental.Facility{
		Id:           domain.NewId[facilityrental.Facility](dto.FacilityID),
		FacilityType: facilityType,
	}

	validity := facilityrental.RentalValidity{
		ToDate: dto.ExpiresAt,
	}

	// Check if this is a boat facility (has boat info)
	if dto.BoatID != nil && dto.BoatName != nil && dto.LengthMeters != nil && dto.WidthMeters != nil {
		boatInfo := facilityrental.BoatInfo{
			Name:          *dto.BoatName,
			LengthMeters:  *dto.LengthMeters,
			WidthMeters:   *dto.WidthMeters,
			InsuranceInfo: facilityrental.NoBoatInsurance{},
		}

		return facilityrental.RentedFacilityWithBoat{
			Id:       domain.NewId[facilityrental.RentedFacility](dto.RentedFacilityID),
			MemberId: domain.NewId[membership.Member](0), // Will be filled from query param
			Facility: facility.Id,
			Validity: validity,
			Payment:  nil, // Payment info not included in this query
			BoatInfo: boatInfo,
		}
	}

	// Simple facility without boat
	return facilityrental.SimpleRentedFacility{
		Id:       domain.NewId[facilityrental.RentedFacility](dto.RentedFacilityID),
		MemberId: domain.NewId[membership.Member](0), // Will be filled from query param
		Facility: facility.Id,
		Validity: validity,
		Payment:  nil, // Payment info not included in this query
	}
}
