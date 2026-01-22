package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/club"
	facilityrental "github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/reports"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/persistence"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/presentation"
	infrareports "github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/reports"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

var (
	memberService      *membership.MemberManagementService
	rentalService      *facilityrental.RentalManagementService
	paymentService     *payment.PaymentManagementService
	waitingListService *facilityrental.WaitingListManagementService
	reportService      *reports.ReportService
	facilityRepo       facilityrental.FacilityRepository
	seasonRepo         club.SeasonRepository
)

func InitializeServices(database *sql.DB) {
	var memberRepository = persistence.NewSQLMemberRepository(database)
	memberService = membership.NewMemberManagementService(memberRepository)
	facilityRepo = persistence.NewSQLFacilityRepository(database)
	waitingListRepo := persistence.NewSQLWaitingListRepository(database)
	rentalService = facilityrental.NewRentalManagementService(facilityRepo, waitingListRepo)
	paymentRepo := persistence.NewSQLPaymentRepository(database)
	paymentService = payment.NewPaymentManagementService(paymentRepo)
	waitingListService = facilityrental.NewWaitingListManagementService(waitingListRepo)
	seasonRepo = persistence.NewSQLSeasonRepository(database)
	pdfGenerator := infrareports.NewChromeDPPDFGenerator()
	reportService = reports.NewReportService(pdfGenerator)
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

		// Convert boat info from presentation to domain model if provided
		var boatInfo *facilityrental.BoatInfo = nil
		if req.BoatInfo != nil {
			// Validate boat info
			if req.BoatInfo.Name == "" {
				presentation.WriteError(w, http.StatusBadRequest, "boat name is required when boatInfo is provided")
				return
			}
			if req.BoatInfo.LengthMeters <= 0 {
				presentation.WriteError(w, http.StatusBadRequest, "boat length must be greater than 0")
				return
			}
			if req.BoatInfo.WidthMeters <= 0 {
				presentation.WriteError(w, http.StatusBadRequest, "boat width must be greater than 0")
				return
			}
			if len(req.BoatInfo.Insurances) == 0 {
				presentation.WriteError(w, http.StatusBadRequest, "at least one insurance is required for boat")
				return
			}

			// Use the first insurance (frontend sends array with one item)
			insurance := req.BoatInfo.Insurances[0]
			if insurance.Provider == "" || insurance.Number == "" || insurance.ExpiresAt == "" {
				presentation.WriteError(w, http.StatusBadRequest, "insurance provider, number, and expiration date are required")
				return
			}

			boatInfo = &facilityrental.BoatInfo{
				Name:         req.BoatInfo.Name,
				LengthMeters: req.BoatInfo.LengthMeters,
				WidthMeters:  req.BoatInfo.WidthMeters,
				InsuranceInfo: facilityrental.BoatInsurance{
					ProviderName:   insurance.Provider,
					PolicyNumber:   insurance.Number,
					ExpirationDate: insurance.ExpiresAt,
				},
			}
		}

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

func RentedFacilityByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1.0/facilities/rented/")
	if idStr == "" {
		presentation.WriteError(w, http.StatusBadRequest, "missing rented facility id")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid rented facility id format")
		return
	}

	rentedFacilityId := domain.Id[facilityrental.RentedFacility]{Value: id}

	switch r.Method {
	case http.MethodDelete:
		if rentalService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		// Free the facility
		result := rentalService.FreeFacility(rentedFacilityId)
		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusNotFound, result.Error().Error())
			return
		}

		presentation.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})

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

func WaitingListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if waitingListService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		// Get facility_type_id from query parameter
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

		facilityTypeId := domain.Id[facilityrental.FacilityType]{Value: facilityTypeID}
		result := waitingListService.GetWaitingList(facilityTypeId)

		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusInternalServerError, result.Error().Error())
			return
		}

		waitingList := presentation.ConvertWaitingListToPresentation(result.Value())
		presentation.WriteJSON(w, http.StatusOK, waitingList)

	case http.MethodPost:
		if waitingListService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		var req presentation.AddToWaitingListRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}

		// Validate required fields
		if req.MemberId == 0 {
			presentation.WriteError(w, http.StatusBadRequest, "memberId is required")
			return
		}
		if req.FacilityTypeId == 0 {
			presentation.WriteError(w, http.StatusBadRequest, "facilityTypeId is required")
			return
		}

		memberId := domain.Id[membership.Member]{Value: req.MemberId}
		facilityTypeId := domain.Id[facilityrental.FacilityType]{Value: req.FacilityTypeId}

		result := waitingListService.AddToWaitingList(memberId, facilityTypeId, req.Notes)

		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusInternalServerError, result.Error().Error())
			return
		}

		entry := presentation.ConvertWaitingListEntryToPresentation(result.Value())
		presentation.WriteJSON(w, http.StatusCreated, entry)

	case http.MethodDelete:
		if waitingListService == nil {
			presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
			return
		}

		// Get member_id and facility_type_id from query parameters
		memberIDStr := r.URL.Query().Get("member_id")
		facilityTypeIDStr := r.URL.Query().Get("facility_type_id")

		if memberIDStr == "" || facilityTypeIDStr == "" {
			presentation.WriteError(w, http.StatusBadRequest, "member_id and facility_type_id are required")
			return
		}

		memberID, err := strconv.ParseInt(memberIDStr, 10, 64)
		if err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "invalid member_id")
			return
		}

		facilityTypeID, err := strconv.ParseInt(facilityTypeIDStr, 10, 64)
		if err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "invalid facility_type_id")
			return
		}

		memberId := domain.Id[membership.Member]{Value: memberID}
		facilityTypeId := domain.Id[facilityrental.FacilityType]{Value: facilityTypeID}

		result := waitingListService.RemoveFromWaitingListByMemberAndType(memberId, facilityTypeId)

		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusNotFound, result.Error().Error())
			return
		}

		entry := presentation.ConvertWaitingListEntryToPresentation(result.Value())
		presentation.WriteJSON(w, http.StatusOK, entry)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func SuggestedPriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if rentalService == nil {
		presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
		return
	}

	// Get facility_type_id from query parameter
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

	// Get member_id from query parameter
	memberIDStr := r.URL.Query().Get("member_id")
	if memberIDStr == "" {
		presentation.WriteError(w, http.StatusBadRequest, "member_id is required")
		return
	}

	memberID, err := strconv.ParseInt(memberIDStr, 10, 64)
	if err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid member_id")
		return
	}

	// Get season from query parameter (optional, defaults to 0)
	seasonID := int64(0)
	seasonStr := r.URL.Query().Get("season")
	if seasonStr != "" {
		parsedSeason, err := strconv.ParseInt(seasonStr, 10, 64)
		if err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "invalid season")
			return
		}
		seasonID = parsedSeason
	}

	// Fetch facility type from catalog to get base price
	catalog := facilityRepo.GetFacilitiesCatalog()
	var facilityType *facilityrental.FacilityType
	for _, ft := range catalog {
		if ft.Id.Value == facilityTypeID {
			facilityType = &ft
			break
		}
	}

	if facilityType == nil {
		presentation.WriteError(w, http.StatusNotFound, "facility type not found")
		return
	}

	basePrice := facilityType.SuggestedPrice

	facilityTypeId := domain.Id[facilityrental.FacilityType]{Value: facilityTypeID}
	memberId := domain.Id[membership.User]{Value: memberID}

	// Calculate suggested price with discounts
	suggestedPrice := rentalService.GetSuggestedPriceForMember(
		facilityTypeId,
		basePrice,
		memberId,
		seasonID,
	)

	// Get applicable discount rules for informational purposes
	applicableDiscounts := rentalService.GetApplicableDiscountsForMember(
		facilityTypeId,
		memberId,
		seasonID,
	)

	hasSpecialPrice := suggestedPrice < basePrice
	savingsAmount := 0.0
	if hasSpecialPrice {
		savingsAmount = basePrice - suggestedPrice
	}

	response := map[string]any{
		"suggestedPrice":  suggestedPrice,
		"basePrice":       basePrice,
		"savingsAmount":   savingsAmount,
		"hasSpecialPrice": hasSpecialPrice,
		"applicableRules": len(applicableDiscounts),
	}

	presentation.WriteJSON(w, http.StatusOK, response)
}

// MemberListPDFHandler generates a PDF report with all members
func MemberListPDFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if reportService == nil {
		presentation.WriteError(w, http.StatusInternalServerError, "report service not initialized")
		return
	}

	if memberService == nil {
		presentation.WriteError(w, http.StatusInternalServerError, "member service not initialized")
		return
	}

	// Get season from query parameter
	seasonStr := r.URL.Query().Get("season")
	if seasonStr == "" {
		presentation.WriteError(w, http.StatusBadRequest, "missing season query parameter")
		return
	}

	seasonId, err := strconv.ParseInt(seasonStr, 10, 64)
	if err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid season format")
		return
	}

	// Get all members for the season
	membersResult := memberService.GetListOfMembersBySeason(seasonId)
	if !membersResult.IsSuccess() {
		presentation.WriteError(w, http.StatusInternalServerError, "failed to get members: "+membersResult.Error().Error())
		return
	}

	members := membersResult.Value()

	// Convert to report format
	memberSummaries := make([]reports.MemberSummary, len(members))
	for i, member := range members {
		memberSummaries[i] = reports.MemberSummary{
			ID:                  member.User.Id.Value,
			FirstName:           member.User.FirstName,
			LastName:            member.User.LastName,
			Email:               member.User.Email.Value,
			BirthDate:           member.User.BirthDate.Format("02/01/2006"),
			MembershipNumber:    member.Membership.Number,
			MembershipStatus:    string(member.Membership.Status.GetStatus()),
			MembershipPaid:      member.Membership.Payment.GetStatus() == payment.Paid,
			HasUnpaidFacilities: member.HasUnpaidFacilities,
		}
	}

	// Get season code
	seasonResult := seasonRepo.GetSeasonById(seasonId)
	if !seasonResult.IsSuccess() {
		presentation.WriteError(w, http.StatusInternalServerError, "failed to get season: "+seasonResult.Error().Error())
		return
	}
	seasonCode := seasonResult.Value().GetCode()

	// Generate PDF
	pdfBuffer, err := reportService.GenerateMemberListReport(memberSummaries, seasonCode)
	if err != nil {
		presentation.WriteError(w, http.StatusInternalServerError, "failed to generate PDF: "+err.Error())
		return
	}

	// Set headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=lista_soci.pdf")
	w.Header().Set("Content-Length", strconv.Itoa(pdfBuffer.Len()))

	// Write PDF to response
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBuffer.Bytes())
}

// MemberDetailPDFHandler generates a PDF report with member details and facilities
func MemberDetailPDFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if reportService == nil {
		presentation.WriteError(w, http.StatusInternalServerError, "report service not initialized")
		return
	}

	if memberService == nil || rentalService == nil {
		presentation.WriteError(w, http.StatusInternalServerError, "service not initialized")
		return
	}

	// Get member ID from path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1.0/reports/members/")
	idStr = strings.TrimSuffix(idStr, "/pdf")
	if idStr == "" {
		presentation.WriteError(w, http.StatusBadRequest, "missing member id")
		return
	}

	memberId64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid member id format")
		return
	}

	memberId := domain.Id[membership.Member]{Value: memberId64}

	// Get season from query parameter
	seasonStr := r.URL.Query().Get("season")
	if seasonStr == "" {
		presentation.WriteError(w, http.StatusBadRequest, "missing season query parameter")
		return
	}

	seasonId, err := strconv.ParseInt(seasonStr, 10, 64)
	if err != nil {
		presentation.WriteError(w, http.StatusBadRequest, "invalid season format")
		return
	}

	// Get member details
	memberResult := memberService.GetMemberById(memberId, seasonId)
	if !memberResult.IsSuccess() {
		presentation.WriteError(w, http.StatusNotFound, "member not found")
		return
	}

	memberDetails := memberResult.Value()

	// Convert member to report format
	memberDetail := reports.MemberDetail{
		ID:        memberDetails.User.Id.Value,
		FirstName: memberDetails.User.FirstName,
		LastName:  memberDetails.User.LastName,
		Email:     memberDetails.User.Email.Value,
		BirthDate: memberDetails.User.BirthDate.Format("02/01/2006"),
	}

	// Add phone numbers
	for _, phone := range memberDetails.User.PhoneNumbers {
		prefix := ""
		if phone.Prefix != nil {
			prefix = *phone.Prefix
		}
		memberDetail.PhoneNumbers = append(memberDetail.PhoneNumbers, reports.PhoneNumber{
			Prefix: prefix,
			Number: phone.Number,
		})
	}

	// Add addresses
	for _, addr := range memberDetails.User.Addresses {
		memberDetail.Addresses = append(memberDetail.Addresses, reports.Address{
			Country:      addr.Country,
			City:         addr.City,
			Street:       addr.Street,
			StreetNumber: addr.Number,
			ZipCode:      addr.ZipCode,
		})
	}

	// Add memberships
	for _, ms := range memberDetails.Memberships {
		memberDetail.Memberships = append(memberDetail.Memberships, reports.Membership{
			ID:        ms.Id.Value,
			Number:    ms.Number,
			Status:    string(ms.Status.GetStatus()),
			ValidFrom: ms.Status.GetValidFromDate().Format("02/01/2006"),
			ExpiresAt: ms.Status.GetValidUntilDate().Format("02/01/2006"),
			Price:     ms.Price,
			Paid:      ms.Payment.GetStatus() == payment.Paid,
		})
	}

	// Get rented facilities
	userIdForFacilities := domain.Id[membership.User]{Value: memberId64}
	rentedFacilities := rentalService.GetFacilitiesRentedByMember(userIdForFacilities, seasonId)

	// Convert to report format
	facilityRentals := make([]reports.FacilityRental, len(rentedFacilities))
	for i, rf := range rentedFacilities {
		facility := rf.GetFacility()
		boatName := ""
		if rf.GetType() == facilityrental.BoatFacility {
			if boatFacility, ok := rf.(facilityrental.RentedFacilityWithBoat); ok {
				boatName = boatFacility.BoatInfo.Name
			}
		}

		facilityRentals[i] = reports.FacilityRental{
			ID:                      rf.GetId().Value,
			FacilityIdentifier:      facility.Identifier,
			FacilityName:            string(facility.FacilityType.FacilityName),
			FacilityTypeDescription: facility.FacilityType.Description,
			RentedAt:                rf.GetValidity().FromDate.Format("02/01/2006"),
			ExpiresAt:               rf.GetValidity().ToDate.Format("02/01/2006"),
			Price:                   rf.GetPrice(),
			Paid:                    rf.GetPayment().GetStatus() == payment.Paid,
			BoatName:                boatName,
		}
	}

	// Get season code
	seasonResult := seasonRepo.GetSeasonById(seasonId)
	if !seasonResult.IsSuccess() {
		presentation.WriteError(w, http.StatusInternalServerError, "failed to get season: "+seasonResult.Error().Error())
		return
	}
	seasonCode := seasonResult.Value().GetCode()

	// Generate PDF
	pdfBuffer, err := reportService.GenerateMemberDetailReport(memberDetail, facilityRentals, seasonCode)
	if err != nil {
		presentation.WriteError(w, http.StatusInternalServerError, "failed to generate PDF: "+err.Error())
		return
	}

	// Set headers for PDF download
	filename := "dettaglio_socio_" + memberDetails.User.LastName + "_" + memberDetails.User.FirstName + ".pdf"
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Length", strconv.Itoa(pdfBuffer.Len()))

	// Write PDF to response
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBuffer.Bytes())
}
