package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/persistence"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/presentation"
)

var (
	memberService *membership.MemberManagementService
	facilityRepo  *persistence.SQLFacilityRepository
)

func InitializeServices(db *sql.DB) {
	memberRepository := persistence.NewSQLMemberRepository(db)
	memberService = membership.NewMemberManagementService(memberRepository)
	facilityRepo = persistence.NewSQLFacilityRepository(db)
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

		result := memberService.GetUpdatedListOfMembers()
		if !result.IsSuccess() {
			presentation.WriteError(w, http.StatusInternalServerError, result.Error().Error())
			return
		}

		members := presentation.ConvertMembersToPresentation(result.Value())
		presentation.WriteJSON(w, http.StatusOK, members)

	case http.MethodPost:
		var m presentation.Member
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			presentation.WriteError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		presentation.WriteJSON(w, http.StatusCreated, m)

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

		memberId := domain.Id[membership.Member]{Value: id}
		result := memberService.GetMemberById(memberId)

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

	// Get DTOs from repository
	dtos, err := facilityRepo.GetRentedFacilityDTOs(memberId)
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
			Payment:                 payment,
			BoatInfo:                boatInfo,
		}
	}

	presentation.WriteJSON(w, http.StatusOK, rentedFacilities)
}
