package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	facilityrental "github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/persistence"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/presentation"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

var (
	memberService *membership.MemberManagementService
	rentalService *facilityrental.RentalManagementService
)

func InitializeServices(db *sql.DB) {
	var memberRepository = persistence.NewSQLMemberRepository(db)
	memberService = membership.NewMemberManagementService(memberRepository)
	facilityRepo := persistence.NewSQLFacilityRepository(db)
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
		seasonId, err := strconv.ParseInt(season, 10, 64)
		if err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "missing season query parameter")
			return
		}

		memberId := domain.Id[membership.Member]{Value: id}
		result := memberService.GetMemberById(memberId, seasonId)

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
	switch r.Method {
	case http.MethodGet:
		if rentalService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		// Get member_id from query parameter
		memberIdStr := r.URL.Query().Get("member_id")
		if memberIdStr == "" {
			presentation.WriteError(w, http.StatusBadRequest, "missing member_id query parameter")
			return
		}

		memberIdValue, err := strconv.ParseInt(memberIdStr, 10, 64)
		if err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "invalid member_id format")
			return
		}
		memberId := domain.NewId[membership.User](memberIdValue)

		// Get season from query parameter
		season := r.URL.Query().Get("season")
		seasonId, err := strconv.ParseInt(season, 10, 64)
		if err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "missing season query parameter")
			return
		}

		// Get DTOs from repository
		rentedFacilities := rentalService.GetFacilitiesRentedByMember(memberId, seasonId)

		// Convert DTOs to presentation models
		rentedFacilitiesDTOs := make([]presentation.RentedFacility, len(rentedFacilities))
		for i, rentedFacility := range rentedFacilities {
			rentedFacilitiesDTOs[i] = presentation.ConvertRentedFacilityToPresentation(rentedFacility)
		}

		presentation.WriteJSON(w, http.StatusOK, rentedFacilitiesDTOs)

	case http.MethodPost:
		if rentalService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		var req presentation.RentFacilityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}

		// Validate required fields
		if req.FacilityId == 0 {
			presentation.WriteError(w, http.StatusBadRequest, "facilityId is required")
			return
		}
		if req.MemberId == 0 {
			presentation.WriteError(w, http.StatusBadRequest, "memberId is required")
			return
		}
		if req.SeasonId == 0 {
			presentation.WriteError(w, http.StatusBadRequest, "seasonId is required")
			return
		}
		if req.RentedAt == "" {
			presentation.WriteError(w, http.StatusBadRequest, "seasonStartsAt is required")
			return
		}
		if req.ExpiresAt == "" {
			presentation.WriteError(w, http.StatusBadRequest, "seasonEndsAt is required")
			return
		}
		if req.Price < 0 {
			presentation.WriteError(w, http.StatusBadRequest, "price must be greater than 0")
			return
		}

		rentedAt, err := time.Parse("2006-01-02", req.RentedAt)
		if err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "rentedAt date format not valid")
			return
		}
		expiresAt, err := time.Parse("2006-01-02", req.ExpiresAt)
		if err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "expiresAt date format not valid")
			return
		}

		memberId := domain.Id[membership.User]{Value: req.MemberId}
		facilityId := domain.Id[facilityrental.Facility]{Value: req.FacilityId}
		var boatInfo *facilityrental.BoatInfo = nil

		// Rent facility
		result := rentalService.RentService(
			facilityId,
			memberId,
			req.SeasonId,
			facilityrental.RentalValidity{
				FromDate: rentedAt,
				ToDate:   expiresAt,
			},
			req.Price,
			boatInfo,
		)
		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusInternalServerError, result.Error().Error())
			return
		}

		// Convert to presentation and return
		rentedFacility := presentation.ConvertRentedFacilityToPresentation(result.Value())
		presentation.WriteJSON(w, http.StatusCreated, rentedFacility)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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
	if req.Price < 0 {
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
