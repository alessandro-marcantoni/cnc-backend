package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	facilityrental "github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

//go:embed queries/get_rented_facilities_by_member.sql
var getRentedFacilitiesByMemberQuery string

//go:embed queries/get_facilities_catalog.sql
var getFacilitiesCatalogQuery string

//go:embed queries/get_facilities_by_type.sql
var getFacilitiesByTypeQuery string

//go:embed queries/insert_rented_facility.sql
var insertRentedFacilityQuery string

//go:embed queries/get_facility_pricing_rules.sql
var getFacilityPricingRulesQuery string

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

func (r *SQLFacilityRepository) GetFacilitiesByType(facilityTypeId domain.Id[facilityrental.FacilityType], seasonId int64) []facilityrental.FacilityWithStatus {
	rows, err := r.db.Query(getFacilitiesByTypeQuery, facilityTypeId.Value, seasonId)
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
		var rentedByMemberId sql.NullInt64
		var rentedByMemberFirstName sql.NullString
		var rentedByMemberLastName sql.NullString

		err := rows.Scan(
			&id,
			&facilityTypeId,
			&identifier,
			&facilityTypeName,
			&facilityTypeDescription,
			&suggestedPrice,
			&isRented,
			&expiresAt,
			&rentedByMemberId,
			&rentedByMemberFirstName,
			&rentedByMemberLastName,
		)
		if err != nil {
			continue
		}

		var expiresAtPtr *time.Time
		if expiresAt.Valid {
			expiresAtPtr = &expiresAt.Time
		}

		var memberIdPtr *int64
		if rentedByMemberId.Valid {
			memberIdPtr = &rentedByMemberId.Int64
		}

		var firstNamePtr *string
		if rentedByMemberFirstName.Valid {
			firstNamePtr = &rentedByMemberFirstName.String
		}

		var lastNamePtr *string
		if rentedByMemberLastName.Valid {
			lastNamePtr = &rentedByMemberLastName.String
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
			RentedByMemberId:        memberIdPtr,
			RentedByMemberFirstName: firstNamePtr,
			RentedByMemberLastName:  lastNamePtr,
		}
		facilities = append(facilities, facility)
	}

	return facilities
}

func (r *SQLFacilityRepository) GetAvailableFacilities(serviceType facilityrental.FacilityName) []facilityrental.Facility {
	// TODO: Implement this method
	return []facilityrental.Facility{}
}

func (r *SQLFacilityRepository) GetFacilitiesRentedByMember(memberId domain.Id[membership.User], season int64) []facilityrental.RentedFacility {
	rows, err := r.db.Query(getRentedFacilitiesByMemberQuery, memberId.Value, season)
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
			&dto.Price,
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

		rentedFacility := ConvertDTOToRentedFacility(dto)
		rentedFacilities = append(rentedFacilities, rentedFacility)
	}

	return rentedFacilities
}

func (r *SQLFacilityRepository) RentFacility(
	memberId domain.Id[membership.User],
	facilityId domain.Id[facilityrental.Facility],
	season int64,
	price float64,
	boatInfo *facilityrental.BoatInfo,
) result.Result[facilityrental.RentedFacility] {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result.Err[facilityrental.RentedFacility](errors.RepositoryError{Description: "failed to begin transaction: " + err.Error()})
	}
	defer tx.Rollback()

	// Insert facility rental
	var rentedFacilityId int64
	err = tx.QueryRowContext(ctx, insertRentedFacilityQuery,
		facilityId.Value,
		memberId.Value,
		season,
		price,
	).Scan(&rentedFacilityId)
	if err != nil {
		return result.Err[facilityrental.RentedFacility](errors.RepositoryError{Description: "failed to insert facility rental: " + err.Error()})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return result.Err[facilityrental.RentedFacility](errors.RepositoryError{Description: "failed to commit transaction: " + err.Error()})
	}

	rentedFacilities := r.GetFacilitiesRentedByMember(memberId, season)
	for _, rentedFacility := range rentedFacilities {
		if rentedFacility.GetId().Value == rentedFacilityId {
			return result.Ok(rentedFacility)
		}
	}

	return result.Err[facilityrental.RentedFacility](errors.RepositoryError{Description: "failed to retrieve inserted facility id"})
}

func (r *SQLFacilityRepository) GetPricingRules() []facilityrental.PricingRule {
	rows, err := r.db.Query(getFacilityPricingRulesQuery)
	if err != nil {
		return []facilityrental.PricingRule{}
	}
	defer rows.Close()

	var pricingRules []facilityrental.PricingRule

	for rows.Next() {
		var id int64
		var facilityTypeId int64
		var requiredFacilityTypeId int64
		var specialPrice float64
		var currency string
		var description sql.NullString
		var active bool
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(
			&id,
			&facilityTypeId,
			&requiredFacilityTypeId,
			&specialPrice,
			&currency,
			&description,
			&active,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			continue
		}

		pricingRule := facilityrental.PricingRule{
			Id:                     domain.Id[facilityrental.PricingRule]{Value: id},
			FacilityTypeId:         domain.Id[facilityrental.FacilityType]{Value: facilityTypeId},
			RequiredFacilityTypeId: domain.Id[facilityrental.FacilityType]{Value: requiredFacilityTypeId},
			SpecialPrice:           specialPrice,
			Description:            description.String,
			Active:                 active,
		}
		pricingRules = append(pricingRules, pricingRule)
	}

	return pricingRules
}
