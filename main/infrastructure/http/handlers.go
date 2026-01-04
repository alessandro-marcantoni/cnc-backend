package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/presentation"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	presentation.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func MembersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		members := []presentation.Member{
			{ID: "1", Name: "Mario Rossi", Email: "mario@example.com"},
		}
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
	id := strings.TrimPrefix(r.URL.Path, "/api/v1.0/members/")
	if id == "" {
		presentation.WriteError(w, http.StatusBadRequest, "missing id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		presentation.WriteJSON(w, http.StatusOK, presentation.Member{
			ID:    id,
			Name:  "Mario Bianchi",
			Email: "mario@example.com",
		})

	case http.MethodDelete:
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
