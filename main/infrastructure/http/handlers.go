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
)

func InitializeServices(db *sql.DB) {
	memberRepository := persistence.NewSQLMemberRepository(db)
	memberService = membership.NewMemberManagementService(memberRepository)
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

		member := presentation.ConvertMemberToPresentation(result.Value())
		presentation.WriteJSON(w, http.StatusOK, member)

	case http.MethodDelete:
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
