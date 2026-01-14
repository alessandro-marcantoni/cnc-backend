package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	facilityrental "github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/persistence"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/presentation"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

var (
	memberService *membership.MemberManagementService
	facilityRepo  *persistence.SQLFacilityRepository
	rentalService *facilityrental.RentalManagementService
)

func InitializeServices(db *sql.DB) {
	memberRepository := persistence.NewSQLMemberRepository(db)
	memberService = membership.NewMemberManagementService(memberRepository)
	facilityRepo = persistence.NewSQLFacilityRepository(db)
	rentalService = facilityrental.NewRentalManagementService(facilityRepo)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	presentation.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func MembersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if memberService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		var result result.Result[[]membership.Member]
		switch {
		case r.URL.Query().Get("season") != "":
			year, err := strconv.ParseInt(r.URL.Query().Get("season"), 10, 64)
			if err != nil {
				presentation.WriteError(w, http.StatusBadRequest, "invalid season format")
				return
			}
			result = memberService.GetListOfMembersBySeason(year)
		default:
			result = memberService.GetListOfAllMembers()
		}

		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusInternalServerError, result.Error().Error())
			return
		}

		members := presentation.ConvertMembersToPresentation(result.Value())
		presentation.WriteJSON(w, http.StatusOK, members)

	case http.MethodPost:
		if memberService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		var req presentation.CreateMemberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}

		// Convert presentation request to domain data
		data, err := presentation.ConvertCreateMemberRequestToDomain(req)
		if err != nil {
			presentation.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Create the member
		result := memberService.CreateMember(data.User, data.CreateMembership, data.SeasonId, data.Price)
		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusInternalServerError, result.Error().Error())
			return
		}

		// Convert to presentation and return
		memberDetails := presentation.ConvertMemberDetailsToPresentation(result.Value())
		presentation.WriteJSON(w, http.StatusCreated, memberDetails)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func MemberByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1.0/members/")
	if idStr == "" {
		presentation.WriteError(w, http.StatusBadRequest, "missing id")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	switch r.Method {
	case http.MethodGet:
		if memberService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		// Get season from query parameter
		season := r.URL.Query().Get("season")
		if season == "" {
			presentation.WriteError(w, http.StatusBadRequest, "missing season query parameter")
			return
		}

		memberId := domain.Id[membership.Member]{Value: id}
		result := memberService.GetMemberById(memberId, season)

		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusNotFound, "member not found")
			return
		}

		member := presentation.ConvertMemberDetailsToPresentation(result.Value())
		presentation.WriteJSON(w, http.StatusOK, member)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func RentedFacilitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if facilityRepo == nil {
		presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
		return
	}

	// Get member_id from query parameter
	memberIdStr := r.URL.Query().Get("member_id")
	if memberIdStr == "" {
		presentation.WriteError(w, http.StatusBadRequest, "missing member_id query parameter")
		return
	}

	memberId, err := strconv.ParseInt(memberIdStr, 10, 64)
	if err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid member_id format")
		return
	}

	// Get season from query parameter
	season := r.URL.Query().Get("season")
	if season == "" {
		presentation.WriteError(w, http.StatusBadRequest, "missing season query parameter")
		return
	}

	// Get DTOs from repository
	dtos, err := facilityRepo.GetRentedFacilityDTOs(memberId, season)
	if err != nil {
		presentation.WriteError(w, http.StatusInternalServerError, "failed to retrieve rented facilities")
		return
	}

	// Convert DTOs to presentation models
	rentedFacilities := make([]presentation.RentedFacility, len(dtos))
	for i, dto := range dtos {
		var boatInfo *presentation.BoatInfo
		if dto.BoatID != nil && dto.BoatName != nil && dto.LengthMeters != nil && dto.WidthMeters != nil {
			boatInfo = &presentation.BoatInfo{
				Name:         *dto.BoatName,
				LengthMeters: *dto.LengthMeters,
				WidthMeters:  *dto.WidthMeters,
			}
		}

		var payment *presentation.Payment
		if dto.PaymentID != nil && dto.PaymentAmount != nil {
			payment = &presentation.Payment{
				Amount:   *dto.PaymentAmount,
				Currency: "EUR", // Default currency
			}
			if dto.PaymentCurrency != nil {
				payment.Currency = *dto.PaymentCurrency
			}
			if dto.PaymentPaidAt != nil {
				payment.PaidAt = dto.PaymentPaidAt.Format("2006-01-02T15:04:05Z07:00")
			}
			if dto.PaymentMethod != nil {
				payment.PaymentMethod = *dto.PaymentMethod
			}
			if dto.TransactionRef != nil {
				payment.TransactionRef = *dto.TransactionRef
			}
		}

		rentedFacilities[i] = presentation.RentedFacility{
			ID:                      dto.RentedFacilityID,
			FacilityID:              dto.FacilityID,
			FacilityIdentifier:      dto.FacilityIdentifier,
			FacilityName:            dto.FacilityType,
			FacilityTypeDescription: dto.FacilityTypeDesc,
			RentedAt:                dto.RentedAt.Format("2006-01-02T15:04:05Z07:00"),
			ExpiresAt:               dto.ExpiresAt.Format("2006-01-02"),
			Price:                   dto.Price,
			Payment:                 payment,
			BoatInfo:                boatInfo,
		}
	}

	presentation.WriteJSON(w, http.StatusOK, rentedFacilities)
}

func FacilitiesCatalogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if rentalService == nil {
		presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
		return
	}

	facilityTypes := rentalService.GetFacilitiesCatalog()
	presentationFacilityTypes := presentation.ConvertFacilityTypesToPresentation(facilityTypes)
	presentation.WriteJSON(w, http.StatusOK, presentationFacilityTypes)
}

func FacilitiesByTypeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if rentalService == nil {
		presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
		return
	}

	// Get facility type ID from query parameter
	facilityTypeIDStr := r.URL.Query().Get("facility_type_id")
	if facilityTypeIDStr == "" {
		presentation.WriteError(w, http.StatusBadRequest, "facility_type_id is required")
		return
	}

	facilityTypeID, err := strconv.ParseInt(facilityTypeIDStr, 10, 64)
	if err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid facility_type_id")
		return
	}

	facilities := rentalService.GetFacilitiesByType(domain.Id[facilityrental.FacilityType]{Value: facilityTypeID})
	presentationFacilities := presentation.ConvertFacilitiesWithStatusToPresentation(facilities)
	presentation.WriteJSON(w, http.StatusOK, presentationFacilities)
}

func MembershipsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if memberService == nil {
		presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
		return
	}

	var req presentation.AddMembershipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Validate required fields
	if req.MemberId == 0 {
		presentation.WriteError(w, http.StatusBadRequest, "memberId is required")
		return
	}
	if req.SeasonId == 0 {
		presentation.WriteError(w, http.StatusBadRequest, "seasonId is required")
		return
	}
	if req.SeasonStartsAt == "" {
		presentation.WriteError(w, http.StatusBadRequest, "seasonStartsAt is required")
		return
	}
	if req.SeasonEndsAt == "" {
		presentation.WriteError(w, http.StatusBadRequest, "seasonEndsAt is required")
		return
	}
	if req.Price <= 0 {
		presentation.WriteError(w, http.StatusBadRequest, "price must be greater than 0")
		return
	}

	// Add membership period
	memberId := domain.Id[membership.Member]{Value: req.MemberId}
	result := memberService.AddMembership(memberId, req.SeasonId, req.SeasonStartsAt, req.SeasonEndsAt, req.Price)
	if !result.IsSuccess() {
		presentation.WriteError(w, http.StatusInternalServerError, result.Error().Error())
		return
	}

	// Convert to presentation and return
	memberDetails := presentation.ConvertMemberDetailsToPresentation(result.Value())
	presentation.WriteJSON(w, http.StatusCreated, memberDetails)
}
