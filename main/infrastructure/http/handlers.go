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
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/persistence"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/presentation"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

var (
	memberService  *membership.MemberManagementService
	rentalService  *facilityrental.RentalManagementService
	paymentService *payment.PaymentManagementService
)

func InitializeServices(db *sql.DB) {
	var memberRepository = persistence.NewSQLMemberRepository(db)
	memberService = membership.NewMemberManagementService(memberRepository)
	facilityRepo := persistence.NewSQLFacilityRepository(db)
	rentalService = facilityrental.NewRentalManagementService(facilityRepo)
	paymentRepo := persistence.NewSQLPaymentRepository(db)
	paymentService = payment.NewPaymentManagementService(paymentRepo)
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
			seasonId, err := strconv.ParseInt(r.URL.Query().Get("season"), 10, 64)
			if err != nil {
				presentation.WriteError(w, http.StatusBadRequest, "invalid season ID format")
				return
			}
			result = memberService.GetListOfMembersBySeason(seasonId)
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
		if req.Price < 0 {
			presentation.WriteError(w, http.StatusBadRequest, "price must be greater than 0")
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

	// Get season from query parameter
	seasonStr := r.URL.Query().Get("season")
	if seasonStr == "" {
		presentation.WriteError(w, http.StatusBadRequest, "season is required")
		return
	}

	seasonID, err := strconv.ParseInt(seasonStr, 10, 64)
	if err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid season")
		return
	}

	facilities := rentalService.GetFacilitiesByType(domain.Id[facilityrental.FacilityType]{Value: facilityTypeID}, seasonID)
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
	if req.Price < 0 {
		presentation.WriteError(w, http.StatusBadRequest, "price must be greater than 0")
		return
	}

	// Add membership period
	memberId := domain.Id[membership.Member]{Value: req.MemberId}
	result := memberService.AddMembership(memberId, req.SeasonId, req.Price)
	if !result.IsSuccess() {
		presentation.WriteError(w, http.StatusInternalServerError, result.Error().Error())
		return
	}

	// Convert to presentation and return
	memberDetails := presentation.ConvertMemberDetailsToPresentation(result.Value())
	presentation.WriteJSON(w, http.StatusCreated, memberDetails)
}

func PaymentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if paymentService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		var req presentation.CreatePaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}

		// Validate that exactly one of membershipPeriodId or rentedFacilityId is provided
		if (req.MembershipPeriodId == nil && req.RentedFacilityId == nil) ||
			(req.MembershipPeriodId != nil && req.RentedFacilityId != nil) {
			presentation.WriteError(w, http.StatusBadRequest, "exactly one of membershipPeriodId or rentedFacilityId must be provided")
			return
		}

		// Validate required fields
		if req.Amount < 0 {
			presentation.WriteError(w, http.StatusBadRequest, "amount must be greater than or equal to 0")
			return
		}
		if req.Currency == "" {
			presentation.WriteError(w, http.StatusBadRequest, "currency is required")
			return
		}
		if req.PaymentMethod == "" {
			presentation.WriteError(w, http.StatusBadRequest, "paymentMethod is required")
			return
		}

		var result result.Result[int64]
		if req.MembershipPeriodId != nil {
			result = paymentService.CreatePaymentForMembershipPeriod(
				*req.MembershipPeriodId,
				req.Amount,
				req.Currency,
				req.PaymentMethod,
				req.TransactionRef,
			)
		} else {
			result = paymentService.CreatePaymentForRentedFacility(
				*req.RentedFacilityId,
				req.Amount,
				req.Currency,
				req.PaymentMethod,
				req.TransactionRef,
			)
		}

		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusInternalServerError, result.Error().Error())
			return
		}

		response := map[string]int64{"id": result.Value()}
		presentation.WriteJSON(w, http.StatusCreated, response)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func PaymentByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1.0/payments/")
	if idStr == "" {
		presentation.WriteError(w, http.StatusBadRequest, "missing payment id")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid payment id format")
		return
	}

	paymentId := domain.Id[payment.Payment]{Value: id}

	switch r.Method {
	case http.MethodPut:
		if paymentService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		var req presentation.UpdatePaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}

		// Validate required fields
		if req.Amount < 0 {
			presentation.WriteError(w, http.StatusBadRequest, "amount must be greater than or equal to 0")
			return
		}
		if req.Currency == "" {
			presentation.WriteError(w, http.StatusBadRequest, "currency is required")
			return
		}
		if req.PaymentMethod == "" {
			presentation.WriteError(w, http.StatusBadRequest, "paymentMethod is required")
			return
		}

		result := paymentService.UpdatePayment(
			paymentId,
			req.Amount,
			req.Currency,
			req.PaymentMethod,
			req.TransactionRef,
		)

		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusInternalServerError, result.Error().Error())
			return
		}

		presentation.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})

	case http.MethodDelete:
		if paymentService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		result := paymentService.DeletePayment(paymentId)

		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusNotFound, result.Error().Error())
			return
		}

		presentation.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
